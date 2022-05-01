/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"sync"

	"github.com/andypangaribuan/project9/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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
		if err != nil {
			err = errors.WithStack(err)
		}
	}
	return err
}

func (slf *pqInstance) Execute(sqlQuery string, sqlPars ...interface{}) error {
	_, err := slf.getInstance()
	if err != nil {
		return err
	}

	query, pars := normalizeSqlQueryParams(slf.conn, sqlQuery, sqlPars)
	_, err = execute(slf.conn, query, pars...)
	return err
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
