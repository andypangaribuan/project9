/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"fmt"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (*srDb) NewPostgresInstance(host string, port int, dbName, username, password string, schema *string, config *model.DbConfig, autoRebind, unsafeCompatibility bool, applicationName string) abs.DbPostgresInstance {
	connStr := ""
	if schema == nil {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s application_name=%s sslmode=disable", host, port, username, password, dbName, applicationName)
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable", host, port, username, password, dbName, *schema, applicationName)
	}

	instance := &postgresInstance{
		data: &connModel{
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
		instance.data.maxLifeTimeConnection = config.MaxLifeTimeConnection
		instance.data.maxIdleConnection = config.MaxIdleConnection
		instance.data.maxOpenConnection = config.MaxOpenConnection
	}

	return instance
}

func connect(data *connModel) (*sqlx.DB, error) {
	instance, err := sqlx.Connect(data.driverName, data.connStr)
	if err == nil {
		instance.SetConnMaxLifetime(data.maxLifeTimeConnection)
		instance.SetMaxIdleConns(data.maxIdleConnection)
		instance.SetMaxOpenConns(data.maxOpenConnection)
		err = instance.Ping()
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return instance, err
}

func normalizeSqlQueryParams(data *connModel, sqlQuery string, sqlPars []interface{}) (string, []interface{}) {
	query := sqlQuery
	pars := sqlPars

	if len(sqlPars) == 1 {
		switch val := sqlPars[0].(type) {
		case []interface{}:
			pars = val
		}
	}

	if data.driverName == "postgres" && data.autoRebind && len(pars) > 0 {
		query = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	}
	return query, pars
}
