/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/clog"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
	"github.com/andypangaribuan/project9/p9"
)

type Repo[T any] struct {
	DbInstance        abs.DbInstance
	TableName         string
	ColumnNames       string
	InsertColumnNames string // if empty then use columnNames
	ParamSigns        string
	InsertParamSign   string // if empty then use paramSigns
	PrintUnsafeErr    bool
}

func (slf *Repo[T]) SetColumnNames(names string) {
	slf.ColumnNames = slf.normalizeColumnNames(names)
}

func (slf *Repo[T]) SetInsertColumnNames(names string) {
	slf.InsertColumnNames = slf.normalizeColumnNames(names)
}

func (slf *Repo[T]) OnUnsafe(unsafe *model.DbUnsafeSelectError) {
	if unsafe != nil && slf.PrintUnsafeErr {
		fmt.Printf("[%v] db.unsafe.select.error:\nerror: %v\nsql-query: %v\nsql-pars: %v\ntrace: %v\n",
			p9.Conv.Time.ToStrFull(time.Now()),
			f9.TernaryFnB(unsafe.LogMessage == nil, "nil", func() string { return *unsafe.LogMessage }),
			unsafe.SqlQuery,
			unsafe.SqlPars,
			f9.TernaryFnB(unsafe.LogTrace == nil, "nil", func() string { return *unsafe.LogTrace }),
		)
	}
}

func (slf *Repo[T]) GenerateParamSigns(columnNames string) (paramSign string) {
	cn := strings.TrimSpace(columnNames)
	cn = strings.ReplaceAll(cn, " ", "")
	if cn != "" {
		ls := strings.Split(cn, ",")
		for i := range ls {
			if i > 0 {
				paramSign += ","
			}
			paramSign += " ?"
		}

		paramSign = strings.TrimSpace(paramSign)
	}
	return
}

func (slf *Repo[T]) Insert(logc *clog.Instance, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doInsert(nil, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) InsertRID(logc *clog.Instance, sqlPars ...interface{}) (*int64, error) {
	startAt := f9.TimeNow()
	sql, id, err := slf.doInsertRID(nil, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return id, err
}

func (slf *Repo[T]) TxInsert(logc *clog.Instance, tx abs.DbTx, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doInsert(tx, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) TxInsertRID(logc *clog.Instance, tx abs.DbTx, sqlPars ...interface{}) (*int64, error) {
	startAt := f9.TimeNow()
	sql, id, err := slf.doInsertRID(tx, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return id, err
}

func (slf *Repo[T]) Update(logc *clog.Instance, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doUpdate(nil, keyVals, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) TxUpdate(logc *clog.Instance, tx abs.DbTx, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doUpdate(tx, keyVals, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) GetData(logc *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(nil, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) GetDatas(logc *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(nil, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) TxGetData(logc *clog.Instance, tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(tx, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) TxGetDatas(logc *clog.Instance, tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(tx, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) SelectFirst(logc *clog.Instance, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(nil, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) Select(logc *clog.Instance, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(nil, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) TxSelectFirst(logc *clog.Instance, tx abs.DbTx, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(tx, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) TxSelect(logc *clog.Instance, tx abs.DbTx, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(tx, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) Count(logc *clog.Instance, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doCount(nil, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) RawCount(logc *clog.Instance, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doRawCount(nil, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) TxCount(logc *clog.Instance, tx abs.DbTx, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doCount(tx, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) TxRawCount(logc *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doRawCount(tx, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) Delete(logc *clog.Instance, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doDelete(nil, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}
func (slf *Repo[T]) TxDelete(logc *clog.Instance, tx abs.DbTx, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doDelete(tx, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) Execute(logc *clog.Instance, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doExecute(nil, sqlQuery, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) TxExecute(logc *clog.Instance, tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doExecute(tx, sqlQuery, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}
