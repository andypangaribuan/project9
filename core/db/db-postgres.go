/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"errors"
	"sync"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var locking sync.Mutex
var lockingConnRead sync.Mutex

func (slf *pqInstance) getRWInstance() (*srConnection, *sqlx.DB, string, error) {
	locking.Lock()
	defer locking.Unlock()

	if slf.connRW.instance != nil {
		return slf.connRW, slf.connRW.instance, slf.connRW.host, nil
	}

	instance, dbHost, err := getConnection(slf.connRW)
	if err == nil {
		slf.connRW.instance = instance
	}

	return slf.connRW, instance, dbHost, err
}

func (slf *pqInstance) getROInstance(rw_force ...bool) (*srConnection, *sqlx.DB, string, error) {
	if slf.connRO == nil || (len(rw_force) > 0 && rw_force[0]) {
		return slf.getRWInstance()
	}

	lockingConnRead.Lock()
	defer lockingConnRead.Unlock()

	if slf.connRO.instance != nil {
		return slf.connRO, slf.connRO.instance, slf.connRO.host, nil
	}

	instance, dbHost, err := getConnection(slf.connRO)
	if err == nil {
		slf.connRO.instance = instance
	}

	return slf.connRO, instance, dbHost, err
}

func (slf *pqInstance) Ping() error {
	_, instance, _, err := slf.getRWInstance()
	if err == nil {
		err = instance.Ping()
	}

	return err
}

func (slf *pqInstance) PingRead() error {
	_, instance, _, err := slf.getROInstance()
	if err == nil {
		err = instance.Ping()
	}

	return err
}

func (slf *pqInstance) Execute(sqlQuery string, sqlPars ...interface{}) (string, error) {
	conn, _, dbHost, err := slf.getRWInstance()
	if err != nil {
		return dbHost, err
	}

	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	_, dbHost, err = execute(conn, nil, query, pars...)
	printSql(conn, startTime, query, pars...)
	return dbHost, err
}

func (slf *pqInstance) ExecuteRID(sqlQuery string, sqlPars ...interface{}) (*int64, string, error) {
	conn, _, dbHost, err := slf.getRWInstance()
	if err != nil {
		return nil, dbHost, err
	}

	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	rid, dbHost, err := executeRID(conn, nil, query, pars...)
	printSql(conn, startTime, query, pars...)

	return rid, dbHost, err
}

func (slf *pqInstance) TxExecute(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (string, error) {
	query, pars := normalizeSqlQueryParams(slf.connRW, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	_, dbHost, err := execute(slf.connRW, tx, query, pars...)
	printSql(slf.connRW, startTime, query, pars...)
	return dbHost, err
}

func (slf *pqInstance) TxExecuteRID(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (*int64, string, error) {
	query, pars := normalizeSqlQueryParams(slf.connRW, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	rid, dbHost, err := executeRID(slf.connRW, tx, query, pars...)
	printSql(slf.connRW, startTime, query, pars...)
	return rid, dbHost, err
}

func (slf *pqInstance) Select(rw_force bool, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, string, error) {
	conn, instance, dbHost, err := slf.getROInstance(rw_force)
	if err != nil {
		return nil, dbHost, err
	}

	// var conn *srConnection
	// switch {
	// case rw_force:
	// 	conn = slf.connRW
	// case slf.connRO != nil:
	// 	conn = slf.connRO
	// default:
	// 	conn = slf.connRW
	// }

	// conn := f9.Ternary(slf.connRO != nil, slf.connRO, slf.connRW)
	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	err = instance.Select(out, query, pars...)
	printSql(conn, startTime, query, pars...)
	if err != nil {
		if slf.connRW.unsafeCompatibility {
			msg := err.Error()
			// TODO: implement LogTrace
			unsafe := model.DbUnsafeSelectError{
				LogType:    "error",
				SqlQuery:   query,
				SqlPars:    pars,
				LogMessage: &msg,
				LogTrace:   nil,
			}

			err = instance.Unsafe().Select(out, query, pars...)
			return &unsafe, dbHost, err
		}
	}

	return nil, dbHost, err
}

func (slf *pqInstance) TxSelect(tx abs.DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, string, error) {
	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			query, pars := normalizeSqlQueryParams(slf.connRW, sqlQuery, sqlPars)
			startTime := f9.TimeNow()
			err := v.tx.Select(out, query, pars...)
			printSql(slf.connRW, startTime, query, pars...)
			if err != nil {
				if slf.connRW.unsafeCompatibility {
					msg := err.Error()
					// TODO: implement LogTrace
					unsafe := model.DbUnsafeSelectError{
						LogType:    "error",
						SqlQuery:   query,
						SqlPars:    pars,
						LogMessage: &msg,
						LogTrace:   nil,
					}

					err = v.tx.Unsafe().Select(out, query, pars...)
					return &unsafe, tx.Host(), err
				}
			}

			return nil, tx.Host(), err
		}
	}

	return nil, tx.Host(), errors.New("unknown: tx transaction type")
}

func (slf *pqInstance) Get(rw_force bool, out interface{}, sqlQuery string, sqlPars ...interface{}) (string, error) {
	conn, instance, dbHost, err := slf.getROInstance(rw_force)
	if err != nil {
		return dbHost, err
	}

	// var conn *srConnection
	// switch {
	// case rw_force:
	// 	conn = slf.connRW
	// case slf.connRO != nil:
	// 	conn = slf.connRO
	// default:
	// 	conn = slf.connRW
	// }

	// conn := f9.Ternary(slf.connRO != nil, slf.connRO, slf.connRW)
	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	err = instance.Get(out, query, pars...)
	printSql(conn, startTime, query, pars...)
	return dbHost, err
}

func (slf *pqInstance) TxGet(tx abs.DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) (string, error) {
	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			query, pars := normalizeSqlQueryParams(slf.connRW, sqlQuery, sqlPars)
			startTime := f9.TimeNow()
			err := v.tx.Get(out, query, pars...)
			printSql(slf.connRW, startTime, query, pars...)
			return tx.Host(), err
		}
	}
	return tx.Host(), errors.New("unknown: tx transaction type")
}

func (slf *pqInstance) NewTransaction() (abs.DbTx, string, error) {
	conn, _, dbHost, err := slf.getRWInstance()
	if err != nil {
		return nil, dbHost, err
	}

	ins := &pqInstanceTx{
		isCommit:   false,
		isRollback: false,
	}

	tx, err := conn.instance.Beginx()
	if err != nil {
		return ins, dbHost, err
	}

	ins.instance = slf
	ins.tx = tx

	return ins, dbHost, nil
}

func (slf *pqInstance) EmptyTransaction() (abs.DbTx, string) {
	ins := &pqInstanceTx{
		instance:   slf,
		isCommit:   false,
		isRollback: false,
	}

	return ins, slf.connRW.host
}
