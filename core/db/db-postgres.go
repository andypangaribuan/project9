/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"fmt"
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
	"github.com/pkg/errors"
)

var locking sync.Mutex
var lockingConnRead sync.Mutex

const renewMaxCount = 3

func (slf *pqInstance) getRWInstance(renew bool) (*srConnection, *sqlx.DB, string, error) {
	locking.Lock()
	defer locking.Unlock()

	if slf.connRW.instance != nil && !renew {
		return slf.connRW, slf.connRW.instance, slf.connRW.host, nil
	}

	instance, dbHost, err := getConnection(slf.connRW)
	if err == nil {
		slf.connRW.instance = instance
	}

	return slf.connRW, instance, dbHost, err
}

func (slf *pqInstance) getROInstance(renew bool, rw_force ...bool) (*srConnection, *sqlx.DB, string, error) {
	if slf.connRO == nil || (len(rw_force) > 0 && rw_force[0]) {
		return slf.getRWInstance(renew)
	}

	lockingConnRead.Lock()
	defer lockingConnRead.Unlock()

	if slf.connRO.instance != nil && !renew {
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

	if err != nil {
		if strings.Contains(err.Error(), "canceling statement due to conflict with recovery") ||
			strings.Contains(err.Error(), "unexpected message 'E'; expected ReadyForQuery") {
			return true
		}
	}

	return false
}

func (slf *pqInstance) Ping() error {
	var (
		err        error
		renew      = false
		renewCount = 0
		instance   *sqlx.DB
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		_, instance, _, err = slf.getRWInstance(renew)
		if err == nil {
			err = instance.Ping()
		}

		if err != nil {
			renew = slf.renewConnection(err)
			err = errors.WithStack(err)
		}

		if err == nil || !renew {
			break
		}
	}

	return err
}

func (slf *pqInstance) PingRead() error {
	var (
		err        error
		renew      = false
		renewCount = 0
		instance   *sqlx.DB
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		_, instance, _, err = slf.getROInstance(renew)
		if err == nil {
			err = instance.Ping()
		}

		if err != nil {
			renew = slf.renewConnection(err)
			err = errors.WithStack(err)
		}

		if err == nil || !renew {
			break
		}
	}

	return err
}

func (slf *pqInstance) Execute(sqlQuery string, sqlPars ...interface{}) (string, error) {
	var (
		err        error
		renew      = false
		renewCount = 0
		conn       *srConnection
		dbHost     string
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		conn, _, dbHost, err = slf.getRWInstance(renew)
		if err != nil {
			renew = slf.renewConnection(err)
			if renew {
				continue
			}

			return dbHost, err
		}

		query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
		startTime := f9.TimeNow()
		_, dbHost, err = execute(conn, nil, query, pars...)
		printSql(conn, startTime, query, pars...)

		renew = slf.renewConnection(err)
		if err == nil || !renew {
			break
		}
	}

	return dbHost, err
}

func (slf *pqInstance) ExecuteRID(sqlQuery string, sqlPars ...interface{}) (*int64, string, error) {
	var (
		err        error
		renew      = false
		renewCount = 0
		conn       *srConnection
		dbHost     string
		rid        *int64
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		conn, _, dbHost, err = slf.getRWInstance(renew)
		if err != nil {
			renew = slf.renewConnection(err)
			if renew {
				continue
			}

			return nil, dbHost, err
		}

		query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
		startTime := f9.TimeNow()
		rid, dbHost, err = executeRID(conn, nil, query, pars...)
		printSql(conn, startTime, query, pars...)

		renew = slf.renewConnection(err)
		if err == nil || !renew {
			break
		}
	}

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
	var (
		err        error
		renew      = false
		renewCount = 0
		conn       *srConnection
		instance   *sqlx.DB
		dbHost     string
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		conn, instance, dbHost, err = slf.getROInstance(renew, rw_force)
		if err != nil {
			renew = slf.renewConnection(err)
			if renew {
				continue
			}

			return nil, dbHost, err
		}

		query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
		startTime := f9.TimeNow()
		err = instance.Select(out, query, pars...)
		printSql(conn, startTime, query, pars...)

		if err == nil {
			break
		} else {
			renew = slf.renewConnection(err)
			err = errors.WithStack(err)

			if conn.unsafeCompatibility {
				var (
					msg    = err.Error()
					trace  = fmt.Sprintf("%+v", err)
					unsafe = model.DbUnsafeSelectError{
						LogType:    "error",
						SqlQuery:   query,
						SqlPars:    pars,
						LogMessage: &msg,
						LogTrace:   &trace,
					}
				)

				err = instance.Unsafe().Select(out, query, pars...)
				if err != nil {
					renew = slf.renewConnection(err)
					err = errors.WithStack(err)
					if renew {
						continue
					}
				}

				return &unsafe, dbHost, err
			}

			if renew {
				continue
			} else {
				break
			}
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
				err = errors.WithStack(err)

				if slf.connRW.unsafeCompatibility {
					msg := err.Error()
					trace := fmt.Sprintf("%+v", err)

					unsafe := model.DbUnsafeSelectError{
						LogType:    "error",
						SqlQuery:   query,
						SqlPars:    pars,
						LogMessage: &msg,
						LogTrace:   &trace,
					}

					err = v.tx.Unsafe().Select(out, query, pars...)
					if err != nil {
						err = errors.WithStack(err)
					}

					return &unsafe, tx.Host(), err
				}
			}

			return nil, tx.Host(), err
		}
	}

	return nil, tx.Host(), errors.New("unknown: tx transaction type")
}

func (slf *pqInstance) Get(rw_force bool, out interface{}, sqlQuery string, sqlPars ...interface{}) (string, error) {
	var (
		err        error
		renew      = false
		renewCount = 0
		conn       *srConnection
		instance   *sqlx.DB
		dbHost     string
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		conn, instance, dbHost, err = slf.getROInstance(renew, rw_force)
		if err != nil {
			renew = slf.renewConnection(err)
			if renew {
				continue
			}

			return dbHost, err
		}

		query, pars := normalizeSqlQueryParams(conn, sqlQuery, sqlPars)
		startTime := f9.TimeNow()
		err = instance.Get(out, query, pars...)
		printSql(conn, startTime, query, pars...)

		if err == nil {
			break
		} else {
			renew = slf.renewConnection(err)
			err = errors.WithStack(err)
			if !renew {
				break
			}
		}
	}

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

			if err != nil {
				err = errors.WithStack(err)
			}

			return tx.Host(), err
		}
	}
	return tx.Host(), errors.New("unknown: tx transaction type")
}

func (slf *pqInstance) NewTransaction() (abs.DbTx, string, error) {
	var (
		err        error
		renew      = false
		renewCount = 0
		conn       *srConnection
		dbHost     string
		ins        *pqInstanceTx
	)

	for {
		renewCount++
		if renewCount > renewMaxCount {
			break
		}

		conn, _, dbHost, err = slf.getRWInstance(renew)
		if err != nil {
			renew = slf.renewConnection(err)
			if renew {
				continue
			}

			return nil, dbHost, err
		}

		ins = &pqInstanceTx{
			isCommit:   false,
			isRollback: false,
		}

		tx, err := conn.instance.Beginx()
		if err != nil {
			renew = slf.renewConnection(err)
			err = errors.WithStack(err)
			if renew {
				continue
			}

			return ins, dbHost, err
		}

		ins.instance = slf
		ins.tx = tx
		break
	}

	return ins, dbHost, err
}

func (slf *pqInstance) EmptyTransaction() (abs.DbTx, string) {
	ins := &pqInstanceTx{
		instance:   slf,
		isCommit:   false,
		isRollback: false,
	}

	return ins, slf.connRW.host
}

func (slf *pqInstance) renewConnection(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "driver: bad connection") {
		return true
	}

	return false
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
