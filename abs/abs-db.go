/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "github.com/andypangaribuan/project9/model"

type Db interface {
	NewPostgresInstance(host string, port int, dbName, username, password string, schema *string, config *model.DbConfig, autoRebind, unsafeCompatibility bool, applicationName string) DbPostgresInstance
}

type DbInstance interface {
	Ping() error
	Execute(sqlQuery string, sqlPars ...interface{}) error
	Select(out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error)
	Get(out interface{}, sqlQuery string, sqlPars ...interface{}) error
}

type DbPostgresInstance interface {
	DbInstance
}
