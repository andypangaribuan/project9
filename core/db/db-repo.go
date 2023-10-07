/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"strings"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/clog"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/fc"
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
}

func (slf *Repo[T]) SetColumnNames(names string) {
	slf.ColumnNames = slf.normalizeColumnNames(names)
}

func (slf *Repo[T]) SetInsertColumnNames(names string) {
	slf.InsertColumnNames = slf.normalizeColumnNames(names)
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
	sql, dbHost, err := slf.doInsert(nil, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) InsertRID(logc *clog.Instance, sqlPars ...interface{}) (*int64, error) {
	startAt := f9.TimeNow()
	sql, id, dbHost, err := slf.doInsertRID(nil, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return id, err
}

func (slf *Repo[T]) TxInsert(logc *clog.Instance, tx abs.DbTx, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doInsert(tx, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) TxInsertRID(logc *clog.Instance, tx abs.DbTx, sqlPars ...interface{}) (*int64, error) {
	startAt := f9.TimeNow()
	sql, id, dbHost, err := slf.doInsertRID(tx, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return id, err
}

func (slf *Repo[T]) Update(logc *clog.Instance, sr Update) error {
	var (
		ls      = strings.Split(sr.Set, ",")
		keyVals = make(map[string]interface{}, 0)
	)

	for i, v := range ls {
		key := strings.TrimSpace(strings.Split(v, "=")[0])
		keyVals[key] = sr.SetPars[i]
	}

	return slf.DoUpdate(logc, keyVals, sr.Where, sr.WherePars...)
}

func (slf *Repo[T]) UpdateBy(logc *clog.Instance, set, condition string, pars ...interface{}) error {
	var (
		ls            = strings.Split(set, ",")
		keyVals       = make(map[string]interface{}, 0)
		conditionPars = make([]interface{}, 0)
	)

	for i, v := range ls {
		key := strings.TrimSpace(strings.Split(v, "=")[0])
		keyVals[key] = pars[i]
	}

	for i := len(ls); i < len(pars); i++ {
		conditionPars = append(conditionPars, pars[i])
	}

	return slf.DoUpdate(logc, keyVals, condition, conditionPars...)
}

func (slf *Repo[T]) TxUpdateBy(logc *clog.Instance, tx abs.DbTx, set, condition string, pars ...interface{}) error {
	var (
		ls            = strings.Split(set, ",")
		keyVals       = make(map[string]interface{}, 0)
		conditionPars = make([]interface{}, 0)
	)

	for i, v := range ls {
		key := strings.TrimSpace(strings.Split(v, "=")[0])
		keyVals[key] = pars[i]
	}

	for i := len(ls); i < len(pars); i++ {
		conditionPars = append(conditionPars, pars[i])
	}

	return slf.TxUpdate(logc, tx, keyVals, condition, conditionPars...)
}

func (slf *Repo[T]) DoUpdate(logc *clog.Instance, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doUpdate(nil, keyVals, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) TxUpdate(logc *clog.Instance, tx abs.DbTx, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doUpdate(tx, keyVals, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) Fetch(logc *clog.Instance, condition string, args ...interface{}) (*T, error) {
	whereQuery := getWhereQuery(condition, args...)
	endQuery := strings.TrimSpace(getEndQuery(args...) + " LIMIT 1")
	return slf.GetData(logc, whereQuery, endQuery, getPars(args...)...)
}

func (slf *Repo[T]) GetBy(logc *clog.Instance, condition, end string, pars ...interface{}) (*T, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	endQuery := strings.TrimSpace(end + " LIMIT 1")
	return slf.GetData(logc, whereQuery, endQuery, pars...)
}

func (slf *Repo[T]) Fetches(logc *clog.Instance, condition string, args ...interface{}) ([]T, error) {
	whereQuery := getWhereQuery(condition, args...)
	return slf.GetDatas(logc, whereQuery, getEndQuery(args...), getPars(args...)...)
}

func (slf *Repo[T]) GetsBy(logc *clog.Instance, condition, end string, pars ...interface{}) ([]T, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.GetDatas(logc, whereQuery, end, pars...)
}

func (slf *Repo[T]) GetData(logc *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.goGetDatas(nil, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) GetDatas(logc *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.goGetDatas(nil, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return models, err
}

func (slf *Repo[T]) TxGetBy(logc *clog.Instance, tx abs.DbTx, condition, end string, pars ...interface{}) (*T, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	endQuery := strings.TrimSpace(end + " LIMIT 1")
	return slf.TxGetData(logc, tx, whereQuery, endQuery, pars...)
}

func (slf *Repo[T]) TxGetsBy(logc *clog.Instance, tx abs.DbTx, condition, end string, pars ...interface{}) ([]T, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.TxGetDatas(logc, tx, whereQuery, end, pars...)
}

func (slf *Repo[T]) TxGetData(logc *clog.Instance, tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.goGetDatas(tx, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) TxGetDatas(logc *clog.Instance, tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.goGetDatas(tx, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return models, err
}

func (slf *Repo[T]) SelectFirst(logc *clog.Instance, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.doSelect(nil, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) Select(logc *clog.Instance, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.doSelect(nil, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return models, err
}

func (slf *Repo[T]) TxSelectFirst(logc *clog.Instance, tx abs.DbTx, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.doSelect(tx, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) TxSelect(logc *clog.Instance, tx abs.DbTx, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, dbe, err := slf.doSelect(tx, query, args...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return models, err
}

func (slf *Repo[T]) CountBy(logc *clog.Instance, condition, end string, pars ...interface{}) (int, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.Count(logc, whereQuery, end, pars...)
}

func (slf *Repo[T]) Count(logc *clog.Instance, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, dbe, err := slf.doCount(nil, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return count, err
}

func (slf *Repo[T]) RawCount(logc *clog.Instance, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, dbe, err := slf.doRawCount(nil, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return count, err
}

func (slf *Repo[T]) RawInt(logc *clog.Instance, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawInt(nil, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) RawInt64(logc *clog.Instance, query string, pars ...interface{}) (int64, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawInt64(nil, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) RawFloat64(logc *clog.Instance, query string, pars ...interface{}) (float64, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawFloat64(nil, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) RawFCT(logc *clog.Instance, query string, pars ...interface{}) (fc.FCT, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawFCT(nil, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) SumBy(logc *clog.Instance, column, condition, end string, pars ...interface{}) (fc.FCT, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.Sum(logc, column, whereQuery, end, pars...)
}

func (slf *Repo[T]) Sum(logc *clog.Instance, column, whereQuery, endQuery string, wherePars ...interface{}) (fc.FCT, error) {
	startAt := f9.TimeNow()
	sql, val, dbe, err := slf.doSum(nil, column, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return val, err
}

func (slf *Repo[T]) TxCountBy(logc *clog.Instance, tx abs.DbTx, condition, end string, pars ...interface{}) (int, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.TxCount(logc, tx, whereQuery, end, pars...)
}

func (slf *Repo[T]) TxCount(logc *clog.Instance, tx abs.DbTx, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, dbe, err := slf.doCount(tx, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return count, err
}

func (slf *Repo[T]) TxRawCount(logc *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, dbe, err := slf.doRawCount(tx, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return count, err
}

func (slf *Repo[T]) TxRawInt(logc *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawInt(tx, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) TxRawInt64(logc *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (int64, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawInt64(tx, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) TxRawFloat64(logc *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (float64, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawFloat64(tx, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) TxRawFCT(logc *clog.Instance, tx abs.DbTx, query string, pars ...interface{}) (fc.FCT, error) {
	startAt := f9.TimeNow()
	sql, res, dbe, err := slf.doRawFCT(tx, query, pars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return res, err
}

func (slf *Repo[T]) TxSumBy(logc *clog.Instance, tx abs.DbTx, column, condition, end string, pars ...interface{}) (fc.FCT, error) {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.TxSum(logc, tx, column, whereQuery, end, pars...)
}

func (slf *Repo[T]) TxSum(logc *clog.Instance, tx abs.DbTx, column, whereQuery, endQuery string, wherePars ...interface{}) (fc.FCT, error) {
	startAt := f9.TimeNow()
	sql, val, dbe, err := slf.doSum(tx, column, whereQuery, endQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return val, err
}

func (slf *Repo[T]) DeleteBy(logc *clog.Instance, condition string, pars ...interface{}) error {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.Delete(logc, whereQuery, pars...)
}

func (slf *Repo[T]) TxDeleteBy(logc *clog.Instance, tx abs.DbTx, condition string, pars ...interface{}) error {
	whereQuery := f9.Ternary(condition == "", "", "WHERE "+condition)
	return slf.TxDelete(logc, tx, whereQuery, pars...)
}

func (slf *Repo[T]) Delete(logc *clog.Instance, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doDelete(nil, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) TxDelete(logc *clog.Instance, tx abs.DbTx, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doDelete(tx, whereQuery, wherePars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) Execute(logc *clog.Instance, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doExecute(nil, sqlQuery, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) TxExecute(logc *clog.Instance, tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, dbHost, err := slf.doExecute(tx, sqlQuery, sqlPars...)

	if logc != nil && sql != nil {
		execFunc, execPath := slf.getExecFuncPath()
		dbe := &model.DbExec{Host: dbHost}
		go sendDbq(*logc, sql.query, sql.pars, execFunc, execPath, startAt, dbe, err)
	}

	return err
}

func (slf *Repo[T]) getExecFuncPath() (string, string) {
	depth := 2
	execFunc, execPath := p9.Util.GetExecutionInfo(depth)

	for {
		if strings.Contains(execFunc, "github.com/andypangaribuan/project9/core/db.(*") ||
			strings.Contains(execFunc, "/db/repo.(*") {
			depth += 1
			execFunc, execPath = p9.Util.GetExecutionInfo(depth)
			continue
		}

		break
	}

	return execFunc, execPath
}
