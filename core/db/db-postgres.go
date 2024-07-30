/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
	"github.com/andypangaribuan/project9/p9"
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

func (slf *pqInstance) canRetry(rw_force bool, err error) bool {
	if rw_force {
		return false
	}

	if err != nil && strings.Contains(err.Error(), "canceling statement due to conflict with recovery") {
		return true
	}

	return false
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

func (slf *pqInstance) Select(out interface{}, sqlQuery string, sqlPars ...interface{}) error {
	var (
		dbe    = new(model.DbExec)
		loop   = 2
		unsafe *model.DbUnsafeSelectError
		err    error
	)

	for i := 0; i < loop; i++ {
		var (
			rw_force = i == 1
			dbHost   = ""
		)

		unsafe, dbHost, err = slf.DirectSelect(rw_force, out, sqlQuery, sqlPars...)
		if err != nil && !slf.canRetry(rw_force, err) {
			return err
		}

		if dbe.Host != "" {
			dbe.Host += ", "
		}
		dbe.Host += dbHost

		if out != nil {
			kind := reflect.TypeOf(out).Kind()
			if kind == reflect.Pointer {
				val := reflect.ValueOf(out)
				val = val.Elem()
				kind = val.Kind()
				if kind == reflect.Slice {
					length := val.Len()
					if length > 0 {
						break
					}
				}
			}
		}
	}

	slf.OnUnsafe(unsafe)
	return err
}

func (slf *pqInstance) DirectSelect(rw_force bool, out interface{}, sqlQuery string, sqlPars ...interface{}) (*model.DbUnsafeSelectError, string, error) {
	conn, instance, dbHost, err := slf.getROInstance(rw_force)
	if err != nil {
		return nil, dbHost, err
	}

	query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
	startTime := f9.TimeNow()
	err = instance.Select(out, query, pars...)
	printSql(conn, startTime, query, pars...)
	if err != nil {
		if conn.unsafeCompatibility {
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

func (slf *pqInstance) OnUnsafe(unsafe *model.DbUnsafeSelectError) {
	if unsafe != nil {
		if (slf.connRW != nil && slf.connRW.printUnsafeError) || (slf.connRO != nil && slf.connRO.printUnsafeError) {
			log.Printf("[%v] db.unsafe.select.error:\nerror: %v\nsql-query: %v\nsql-pars: %v\ntrace: %v\n",
				p9.Conv.Time.ToStrFull(time.Now()),
				f9.TernaryFnB(unsafe.LogMessage == nil, "nil", func() string { return *unsafe.LogMessage }),
				unsafe.SqlQuery,
				unsafe.SqlPars,
				f9.TernaryFnB(unsafe.LogTrace == nil, "nil", func() string { return *unsafe.LogTrace }),
			)
		}
	}
}
