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
	"github.com/lib/pq"
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

func (slf *Repo[T]) normalizeColumnNames(names string) string {
	lines := strings.Split(names, "\n")
	val := ""
	isFirstLineEmpty := false

	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if idx == 0 {
				isFirstLineEmpty = true
			}
			continue
		}

		if val != "" && isFirstLineEmpty {
			val += "\t"
		}

		val += line
	}

	return val
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

func (slf *Repo[T]) first(ls []T) *T {
	if len(ls) == 0 {
		return nil
	}
	return &ls[0]
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
	sql, err := slf.doInsert(sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) doInsert(sqlPars ...interface{}) (*srRepoSql[T], error) {
	sql := srRepoSql[T]{}.new(`
INSERT INTO ::tableName (
	::insertColumnNames
) VALUES (
	::insertParamSign
)`)

	err := sql.transform(slf)
	if err != nil {
		return sql, err
	}

	return sql, slf.DbInstance.Execute(sql.query, sqlPars)
}

func (slf *Repo[T]) Update(cli *clog.Instance, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doUpdate(keyVals, whereQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) doUpdate(keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) (*srRepoSql[T], error) {
	sql := srRepoSql[T]{}.new(`UPDATE ::tableName SET`)

	idx := -1
	for key, val := range keyVals {
		idx++
		if idx > 0 {
			sql.query += ","
		}
		sql.query += fmt.Sprintf(" %v=?", key)
		sql.pars = append(sql.pars, val)
	}

	whereQuery = strings.TrimSpace(whereQuery)
	wq := strings.ToLower(whereQuery)
	if (len(wq) <= 6) || (len(wq) > 6 && wq[:6] != "where ") {
		whereQuery = fmt.Sprintf("WHERE %v", whereQuery)
	}

	sql.query += fmt.Sprintf(" %v", whereQuery)
	sql.pars = append(sql.pars, wherePars...)

	err := sql.transform(slf)
	if err != nil {
		return sql, err
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars)
}

func (slf *Repo[T]) GetData(cli *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) GetDatas(cli *clog.Instance, whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.goGetDatas(whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) goGetDatas(whereQuery string, endQuery string, wherePars ...interface{}) (*srRepoSql[T], []T, error) {
	var models []T
	sql := srRepoSql[T]{}.new(`SELECT * FROM ::tableName`)

	query := strings.TrimSpace(whereQuery)
	if query != "" {
		wq := strings.ToLower(query)
		if (len(wq) <= 6) || (len(wq) > 6 && wq[:6] != "where ") {
			query = fmt.Sprintf("WHERE %v", query)
		}
	}

	endQuery = strings.TrimSpace(endQuery)
	if endQuery != "" {
		query += " " + endQuery
		query = strings.TrimSpace(query)
	}

	sql.query += " " + query
	sql.pars = wherePars

	err := sql.transform(slf)
	if err != nil {
		return sql, models, err
	}

	unsafe, err := slf.DbInstance.Select(&models, sql.query, sql.pars)
	slf.OnUnsafe(unsafe)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			msg := strings.TrimSpace(e.Message)
			if msg != "" {
				msg += "\n"
			}
			msg += fmt.Sprintf("sql: %v\npars: %v", sql.query, sql.pars)
			e.Message = msg
			err = e
		}
		return sql, nil, err
	}

	return sql, models, nil
}

func (slf *Repo[T]) SelectFirst(cli *clog.Instance, query string, args ...interface{}) (*T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(query, args...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return slf.first(models), err
}

func (slf *Repo[T]) Select(cli *clog.Instance, query string, args ...interface{}) ([]T, error) {
	startAt := f9.TimeNow()
	sql, models, err := slf.doSelect(query, args...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return models, err
}

func (slf *Repo[T]) doSelect(query string, args ...interface{}) (*srRepoSql[T], []T, error) {
	var models []T
	sql := srRepoSql[T]{}.new("")

	sql.query = strings.TrimSpace(query)
	sql.pars = args

	query, pars, err := transformIn(sql.query, sql.pars...)
	if err != nil {
		return sql, models, err
	}

	sql.query = query
	sql.pars = pars

	unsafe, err := slf.DbInstance.Select(&models, sql.query, sql.pars...)
	slf.OnUnsafe(unsafe)

	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			msg := strings.TrimSpace(e.Message)
			if msg != "" {
				msg += "\n"
			}
			msg += fmt.Sprintf("sql: %v\npars: %v", sql.query, sql.pars)
			e.Message = msg
			err = e
		}
		return sql, nil, err
	}

	return sql, models, nil
}

func (slf *Repo[T]) Count(cli *clog.Instance, whereQuery, endQuery string, wherePars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doCount(whereQuery, endQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) RawCount(cli *clog.Instance, query string, pars ...interface{}) (int, error) {
	startAt := f9.TimeNow()
	sql, count, err := slf.doRawCount(query, pars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return count, err
}

func (slf *Repo[T]) doCount(whereQuery, endQuery string, wherePars ...interface{}) (*srRepoSql[T], int, error) {
	sql := srRepoSql[T]{}.new(`SELECT COUNT(1) FROM ::tableName`)

	query := strings.TrimSpace(whereQuery)
	if query != "" {
		wq := strings.ToLower(query)
		if (len(wq) <= 6) || (len(wq) > 6 && wq[:6] != "where ") {
			query = fmt.Sprintf("WHERE %v", query)
		}
	}

	endQuery = strings.TrimSpace(endQuery)
	if endQuery != "" {
		query += " " + endQuery
		query = strings.TrimSpace(query)
	}

	sql.query += " " + query
	sql.pars = wherePars

	err := sql.transform(slf)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = slf.DbInstance.Get(&count, sql.query, sql.pars)
	return sql, count, err
}

func (slf *Repo[T]) doRawCount(query string, pars ...interface{}) (*srRepoSql[T], int, error) {
	sql := srRepoSql[T]{}.new(``)

	sql.query = strings.TrimSpace(query)
	sql.pars = pars

	query, pars, err := transformIn(sql.query, sql.pars...)
	if err != nil {
		return sql, -1, err
	}

	sql.query = query
	sql.pars = pars

	var count int
	err = slf.DbInstance.Get(&count, sql.query, sql.pars)
	return sql, count, err
}

func (slf *Repo[T]) Delete(cli *clog.Instance, whereQuery string, wherePars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doDelete(whereQuery, wherePars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) doDelete(whereQuery string, wherePars ...interface{}) (*srRepoSql[T], error) {
	sql := srRepoSql[T]{}.new(`DELETE FROM ::tableName`)

	query := strings.TrimSpace(whereQuery)
	if query != "" {
		wq := strings.ToLower(query)
		if (len(wq) <= 6) || (len(wq) > 6 && wq[:6] != "where ") {
			query = fmt.Sprintf("WHERE %v", query)
		}
	}

	sql.query += " " + query
	sql.pars = wherePars

	err := sql.transform(slf)
	if err != nil {
		return nil, err
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars)
}

func (slf *Repo[T]) Execute(cli *clog.Instance, sqlQuery string, sqlPars ...interface{}) error {
	startAt := f9.TimeNow()
	sql, err := slf.doExecute(sqlQuery, sqlPars...)

	if cli != nil && sql != nil {
		execFunc, execPath := p9.Util.GetExecutionInfo(1)
		go sendDbq(*cli, sql.query, sql.pars, execFunc, execPath, startAt, err)
	}

	return err
}

func (slf *Repo[T]) doExecute(sqlQuery string, sqlPars ...interface{}) (*srRepoSql[T], error) {
	sql := srRepoSql[T]{}.new(sqlQuery)
	sql.pars = sqlPars

	err := sql.transform(slf)
	if err != nil {
		return nil, err
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars)
}

func sendDbq(cli clog.Instance, sqlQuery string, sqlPars []interface{}, execFunc, execPath string, startAt time.Time, err error) {
	var (
		sqlParsVal *string
		severity   = clog.Info
		errMessage *string
		stackTrace *string
	)

	if val, err := p9.Json.Encode(sqlPars); err == nil {
		sqlParsVal = &val
	}

	if err != nil {
		severity = clog.Error
		errMessage = f9.Ptr(err.Error())
		stackTrace = f9.Ptr(fmt.Sprintf("%+v", err))
	}

	m := clog.SendDbqModel{
		StartAt:    startAt,
		ExecFunc:   &execFunc,
		ExecPath:   &execPath,
		SqlQuery:   sqlQuery,
		SqlPars:    sqlParsVal,
		Error:      errMessage,
		StackTrace: stackTrace,
	}

	clog.SendDbq(0, cli, severity, m, true)
}
