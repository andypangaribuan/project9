/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"errors"
	"sync"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var locking sync.Mutex

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

func (slf *pqInstance) Ping() error {
	instance, err := slf.getInstance()
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
	_, err = execute(slf.conn, nil, query, pars...)
	return err
}

func (slf *pqInstance) ExecuteRID(sqlQuery string, sqlPars ...interface{}) (*int64, error) {
	_, err := slf.getInstance()
	if err != nil {
		return nil, err
	}

	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	return executeRID(slf.conn, nil, query, pars...)
}

func (slf *pqInstance) TxExecute(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) error {
	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	_, err := execute(slf.conn, tx, query, pars...)
	return err
}

func (slf *pqInstance) TxExecuteRID(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (*int64, error) {
	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	return executeRID(slf.conn, tx, query, pars...)
}

func (slf *pqInstance) Select(out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error) {
	instance, err := slf.getInstance()
	if err != nil {
		return nil, err
	}

	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	err = instance.Select(out, query, pars...)
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
			err := v.tx.Select(out, query, pars...)
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
	instance, err := slf.getInstance()
	if err != nil {
		return err
	}

	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	return instance.Get(out, query, pars...)
}

func (slf *pqInstance) TxGet(tx abs.DbTx, out interface{}, sqlQuery string, sqlPars ...interface{}) error {
	if tx != nil {
		switch v := tx.(type) {
		case *pqInstanceTx:
			query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
			return v.tx.Get(out, query, pars...)
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
