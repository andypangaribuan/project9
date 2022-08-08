/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/model"
	"github.com/jmoiron/sqlx"
)

func (*srDb) NewPostgresInstance(host string, port int, dbName, username, password string, schema *string, config *model.DbConfig, autoRebind, unsafeCompatibility bool, applicationName string) abs.DbPostgresInstance {
	connStr := ""
	if schema == nil {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s application_name=%s sslmode=disable", host, port, username, password, dbName, applicationName)
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable", host, port, username, password, dbName, *schema, applicationName)
	}

	instance := &pqInstance{
		conn: &srConnection{
			connStr:               connStr,
			driverName:            "postgres",
			maxLifeTimeConnection: time.Minute * 5,
			maxIdleConnection:     5,
			maxOpenConnection:     100,
			autoRebind:            autoRebind,
			unsafeCompatibility:   unsafeCompatibility,
		},
	}

	if config != nil {
		instance.conn.maxLifeTimeConnection = config.MaxLifeTimeConnection
		instance.conn.maxIdleConnection = config.MaxIdleConnection
		instance.conn.maxOpenConnection = config.MaxOpenConnection
	}

	return instance
}

func getConnection(conn *srConnection) (*sqlx.DB, error) {
	instance, err := sqlx.Connect(conn.driverName, conn.connStr)
	if err == nil {
		instance.SetConnMaxLifetime(conn.maxLifeTimeConnection)
		instance.SetMaxIdleConns(conn.maxIdleConnection)
		instance.SetMaxOpenConns(conn.maxOpenConnection)
		err = instance.Ping()
	}
	return instance, err
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

func execute(conn *srConnection, sqlQuery string, sqlPars ...interface{}) (sql.Result, error) {
	if conn.driverName == "postgres" && conn.autoRebind && len(sqlPars) > 0 {
		sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	}

	stmt, err := conn.instance.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()

	res, err := stmt.Exec(sqlPars...)
	return res, err
}
