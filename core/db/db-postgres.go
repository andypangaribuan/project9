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

func (slf *pqInstance) getInstance() (*sqlx.DB, error) {
	locking.Lock()
	defer locking.Unlock()

	if slf.conn.instance != nil {
		return slf.conn.instance, nil
	}

	instance, err := getConnection(slf.conn)
	if err == nil {
		slf.conn.instance = instance
	}

	return instance, err
}

func (slf *pqInstance) getReadInstance() (*sqlx.DB, error) {
	if slf.connRead == nil {
		return slf.getInstance()
	}

	lockingConnRead.Lock()
	defer lockingConnRead.Unlock()

	if slf.connRead.instance != nil {
		return slf.connRead.instance, nil
	}

	instance, err := getConnection(slf.connRead)
	if err == nil {
		slf.connRead.instance = instance
	}

	return instance, err
}

func (slf *pqInstance) Ping() error {
	instance, err := slf.getInstance()
	if err == nil {
		err = instance.Ping()
	}

	return err
}

func (slf *pqInstance) PingRead() error {
	instance, err := slf.getReadInstance()
	if err == nil {
		err = instance.Ping()
	}

	return err
}

func (slf *pqInstance) Execute(sqlQuery string, sqlPars ...interface{}) error {
	_, err := slf.getInstance()
	if err != nil {
		return err
	}

	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	_, err = execute(slf.conn, nil, query, pars...)
	printSql(slf.conn, startTime, query, pars...)
	return err
}

func (slf *pqInstance) ExecuteRID(sqlQuery string, sqlPars ...interface{}) (*int64, error) {
	_, err := slf.getInstance()
	if err != nil {
		return nil, err
	}

	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	rid, err := executeRID(slf.conn, nil, query, pars...)
	printSql(slf.conn, startTime, query, pars...)

	return rid, err
}

func (slf *pqInstance) TxExecute(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) error {
	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	_, err := execute(slf.conn, tx, query, pars...)
	printSql(slf.conn, startTime, query, pars...)
	return err
}

func (slf *pqInstance) TxExecuteRID(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (*int64, error) {
	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	rid, err := executeRID(slf.conn, tx, query, pars...)
	printSql(slf.conn, startTime, query, pars...)
	return rid, err
}

func (slf *pqInstance) Select(out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error) {
	instance, err := slf.getReadInstance()
	if err != nil {
		return nil, err
	}

	conn := f9.Ternary(slf.connRead != nil, slf.connRead, slf.conn)
	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	err = instance.Select(out, query, pars...)
	printSql(conn, startTime, query, pars...)
	if err != nil {
		if slf.conn.unsafeCompatibility {
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
			return &unsafe, err
		}
	}

	return nil, err
}

func (slf *pqInstance) TxSelect(tx abs.DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error) {
	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
			startTime := f9.TimeNow()
			err := v.tx.Select(out, query, pars...)
			printSql(slf.conn, startTime, query, pars...)
			if err != nil {
				if slf.conn.unsafeCompatibility {
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
					return &unsafe, err
				}
			}

			return nil, err
		}
	}

	return nil, errors.New("unknown: tx transaction type")
}

func (slf *pqInstance) Get(out interface{}, sqlQuery string, sqlPars ...interface{}) error {
	instance, err := slf.getReadInstance()
	if err != nil {
		return err
	}

	conn := f9.Ternary(slf.connRead != nil, slf.connRead, slf.conn)
	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	err = instance.Get(out, query, pars...)
	printSql(conn, startTime, query, pars...)
	return err
}

func (slf *pqInstance) TxGet(tx abs.DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) error {
	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
			startTime := f9.TimeNow()
			err := v.tx.Get(out, query, pars...)
			printSql(slf.conn, startTime, query, pars...)
			return err
		}
	}
	return errors.New("unknown: tx transaction type")
}

func (slf *pqInstance) NewTransaction() (abs.DbTx, error) {
	ins := &pqInstanceTx{
		isCommit:   false,
		isRollback: false,
	}

	tx, err := slf.conn.instance.Beginx()
	if err != nil {
		return ins, err
	}

	ins.instance = slf
	ins.tx = tx

	return ins, nil
}
