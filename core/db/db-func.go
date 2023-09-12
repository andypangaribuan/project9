/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
	"github.com/jmoiron/sqlx"
)

func (*srDb) NewPostgresInstance(host string, port int, dbName, username, password string, schema *string, config *model.DbConfig, autoRebind, unsafeCompatibility bool, applicationName string, printSql bool, printUnsafeError bool) abs.DbPostgresInstance {
	connStr := ""
	if schema == nil {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s application_name=%s sslmode=disable", host, port, username, password, dbName, applicationName)
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable", host, port, username, password, dbName, *schema, applicationName)
	}

	instance := &pqInstance{
		connRW: &srConnection{
			host:                  fmt.Sprintf("%v:%v", host, port),
			connStr:               connStr,
			driverName:            "postgres",
			maxLifeTimeConnection: time.Second * 10,
			maxIdleTimeConnection: time.Second,
			maxIdleConnection:     5,
			maxOpenConnection:     100,
			autoRebind:            autoRebind,
			unsafeCompatibility:   unsafeCompatibility,
			printSql:              printSql,
			printUnsafeError:      printUnsafeError,
		},
	}

	if config != nil {
		instance.connRW.maxLifeTimeConnection = config.MaxLifeTimeConnection
		instance.connRW.maxIdleConnection = config.MaxIdleConnection
		instance.connRW.maxOpenConnection = config.MaxOpenConnection

		if config.MaxIdleTimeConnection > time.Second {
			instance.connRW.maxIdleConnection = config.MaxIdleConnection
		}
	}

	return instance
}

func (*srDb) NewReadWritePostgresInstance(read, write abs.DbPostgresInstance) abs.DbPostgresInstance {
	var (
		connRead  *srConnection
		connWrite *srConnection
	)

	switch r := read.(type) {
	case *pqInstance:
		connRead = r.connRW
	}

	switch w := write.(type) {
	case *pqInstance:
		connWrite = w.connRW
	}

	instance := &pqInstance{
		connRW: connWrite,
		connRO: connRead,
	}

	return instance
}

func getConnection(conn *srConnection) (*sqlx.DB, string, error) {
	instance, err := sqlx.Connect(conn.driverName, conn.connStr)
	if err == nil {
		instance.SetConnMaxLifetime(conn.maxLifeTimeConnection)
		instance.SetConnMaxIdleTime(conn.maxIdleTimeConnection)
		instance.SetMaxIdleConns(conn.maxIdleConnection)
		instance.SetMaxOpenConns(conn.maxOpenConnection)
		err = instance.Ping()
	}

	return instance, conn.host, err
}

func normalizeSqlQueryParams(conn *srConnection, sqlQuery string, sqlPars []interface{}) (string, []interface{}) {
	query := sqlQuery
	pars := sqlPars

	if len(sqlPars) == 1 {
		switch val := sqlPars[0].(type) {
		case []interface{}:
			pars = val
		}
	}

	if conn.driverName == "postgres" && conn.autoRebind && len(pars) > 0 {
		query = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	}
	return query, pars
}

func execute(conn *srConnection, tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (sql.Result, string, error) {
	if conn.driverName == "postgres" && conn.autoRebind && len(sqlPars) > 0 {
		sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	}

	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			stmt, err := v.tx.Prepare(sqlQuery)
			if err != nil {
				return nil, conn.host, err
			}
			defer stmt.Close()

			res, err := stmt.Exec(sqlPars...)
			return res, conn.host, err
		}
	}

	stmt, err := conn.instance.Prepare(sqlQuery)
	if err != nil {
		return nil, conn.host, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(sqlPars...)
	return res, conn.host, err
}

func executeRID(conn *srConnection, tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (*int64, string, error) {
	if conn.driverName == "postgres" && conn.autoRebind && len(sqlPars) > 0 {
		sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	}

	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			stmt, err := v.tx.Prepare(sqlQuery)
			if err != nil {
				return nil, conn.host, err
			}
			defer stmt.Close()

			var id *int64
			err = stmt.QueryRow(sqlPars...).Scan(&id)
			if err != nil {
				return nil, conn.host, err
			}

			return id, conn.host, err
		}
	}

	stmt, err := conn.instance.Prepare(sqlQuery)
	if err != nil {
		return nil, conn.host, err
	}
	defer stmt.Close()

	var id *int64
	err = stmt.QueryRow(sqlPars...).Scan(&id)
	if err != nil {
		return nil, conn.host, err
	}

	return id, conn.host, err
}

func printSql(conn *srConnection, startTime time.Time, sqlQuery string, sqlPars ...interface{}) {
	if conn.printSql {
		durationMs := f9.TimeNow().Sub(startTime).Milliseconds()
		if len(sqlPars) == 0 {
			log.Printf("\nHOST: \"%v\"\nSQL: \"%v\"\nDUR: %vms", conn.host, sqlQuery, durationMs)
		} else {
			log.Printf("\nHOST: \"%v\"\nSQL: \"%v\"\nARG: %v\nDUR: %vms", conn.host, sqlQuery, sqlPars, durationMs)
		}
	}
}
