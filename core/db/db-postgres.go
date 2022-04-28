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

func (slf *postgresInstance) connect() (*sqlx.DB, error) {
	locking.Lock()
	defer locking.Unlock()

	if slf.data.instance != nil {
		return slf.data.instance, nil
	}

	instance, err := connect(slf.data)
	if err == nil {
		slf.data.instance = instance
	}

	return instance, err
}

func (slf *postgresInstance) Ping() error {
	instance, err := slf.connect()
	if err == nil {
		err = instance.Ping()
		if err != nil {
			err = errors.WithStack(err)
		}
	}
	return err
}

func (slf *postgresInstance) Select(out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, error) {
	instance, err := slf.connect()
	if err != nil {
		return nil, err
	}

	query, pars := normalizeSqlQueryParams(slf.data, sqlQuery, sqlPars)
	err = instance.Select(out, query, pars...)
	if err != nil {
		if slf.data.unsafeCompatibility {
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
