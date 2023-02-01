/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "github.com/andypangaribuan/project9/model"

type Db interface {
	NewPostgresInstance(host string, port int, dbName, username, password string, schema *string, config *model.DbConfig, autoRebind, unsafeCompatibility bool, applicationName string) DbPostgresInstance
}

type DbTx interface {
	Commit() error
	Rollback() error
}

type DbInstance interface {
	Ping() error
	Execute(sqlQuery string, sqlPars ...interface{}) error
	ExecuteRID(sqlQuery string, sqlPars ...interface{}) (*int64, error)
	Select(out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error)
	Get(out interface{}, sqlQuery string, sqlPars ...interface{}) error

	NewTransaction() (DbTx, error)

	TxExecute(tx DbTx, sqlQuery string, sqlPars ...interface{}) error
	TxExecuteRID(tx DbTx, sqlQuery string, sqlPars ...interface{}) (*int64, error)
	TxSelect(tx DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error)
	TxGet(tx DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) error
}

type DbPostgresInstance interface {
	DbInstance
}
