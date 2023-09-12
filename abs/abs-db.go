/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "github.com/andypangaribuan/project9/model"

type Db interface {
	NewPostgresInstance(host string, port int, dbName, username, password string, schema *string, config *model.DbConfig, autoRebind, unsafeCompatibility bool, applicationName string, printSql bool, printUnsafeError bool) DbPostgresInstance
	NewReadWritePostgresInstance(read, write DbPostgresInstance) DbPostgresInstance
}

type DbTx interface {
	Commit() error
	Rollback() error
	Host() string
}

type DbInstance interface {
	Ping() error
	Execute(sqlQuery string, sqlPars ...interface{}) (string, error)
	ExecuteRID(sqlQuery string, sqlPars ...interface{}) (*int64, string, error)
	Select(out interface{}, sqlQuery string, sqlPars ...interface{}) error
	DirectSelect(rw_force bool, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, string, error)
	Get(rw_force bool, out interface{}, sqlQuery string, sqlPars ...interface{}) (string, error)

	NewTransaction() (DbTx, string, error)
	EmptyTransaction() (DbTx, string)
	OnUnsafe(unsafe *model.DbUnsafeSelectError)

	TxExecute(tx DbTx, sqlQuery string, sqlPars ...interface{}) (string, error)
	TxExecuteRID(tx DbTx, sqlQuery string, sqlPars ...interface{}) (*int64, string, error)
	TxSelect(tx DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, string, error)
	TxGet(tx DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) (string, error)
}

type DbUpdateHelper interface {
	SetAdd(condition string, param interface{})
	SetAddIfNoNil(condition string, param interface{})
	SetAddIfNotNilOrEmpty(condition string, param interface{})
	Where(condition string, pars ...interface{})

	MapSetAdd(kv map[string]interface{})
	MapSetAddIfNotNilOrEmpty(kv map[string]interface{})

	Get() (condition string, pars []interface{})
}

type DbGetUpdateHelper interface {
	Add(kv map[string]interface{}) DbGetUpdateHelper
	AddIfNotNilOrEmpty(kv map[string]interface{}) DbGetUpdateHelper
	Where(condition string, pars ...interface{}) DbGetUpdateHelper
	Get() (condition string, pars []interface{})
}

type DbPostgresInstance interface {
	DbInstance
}
