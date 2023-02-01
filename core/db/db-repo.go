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

func (slf *Repo[T]) Insert(cli *clog.Instance, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doInsert(nil, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) InsertRID(cli *clog.Instance, sqlPars ...interface{}) (*int64, error) {
	startAt := f9.TimeNow()
	sql, id, err := slf.doInsertRID(nil, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return id, err
}

func (slf *Repo[T]) TxInsert(cli *clog.Instance, tx abs.DbTx, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doInsert(tx, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) TxInsertRID(cli *clog.Instance, tx abs.DbTx, sqlPars ...interface{}) (*int64, error) {
	startAt := f9.TimeNow()
	sql, id, err := slf.doInsertRID(tx, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return id, err
}

func (slf *Repo[T]) Update(cli *clog.Instance, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doUpdate(nil, keyVals, whereQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) TxUpdate(cli *clog.Instance, tx abs.DbTx, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doUpdate(tx, keyVals, whereQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) GetData(cli *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(nil, whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) GetDatas(cli *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(nil, whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) TxGetData(cli *clog.Instance, tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(tx, whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) TxGetDatas(cli *clog.Instance, tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(tx, whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) SelectFirst(cli *clog.Instance, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(nil, query, args...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) Select(cli *clog.Instance, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(nil, query, args...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) TxSelectFirst(cli *clog.Instance, tx abs.DbTx, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(tx, query, args...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) TxSelect(cli *clog.Instance, tx abs.DbTx, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(tx, query, args...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) Count(cli *clog.Instance, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doCount(nil, whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) RawCount(cli *clog.Instance, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doRawCount(nil, query, pars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) TxCount(cli *clog.Instance, tx abs.DbTx, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doCount(tx, whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) TxRawCount(cli *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doRawCount(tx, query, pars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) Delete(cli *clog.Instance, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doDelete(nil, whereQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}
func (slf *Repo[T]) TxDelete(cli *clog.Instance, tx abs.DbTx, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doDelete(tx, whereQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) Execute(cli *clog.Instance, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doExecute(nil, sqlQuery, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) TxExecute(cli *clog.Instance, tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doExecute(tx, sqlQuery, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}
