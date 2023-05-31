/*
 * Copyright (c) 2023.
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
	"github.com/andypangaribuan/project9/fc"
	"github.com/andypangaribuan/project9/model"
	"github.com/andypangaribuan/project9/p9"
	"github.com/lib/pq"
)

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

func (slf *Repo[T]) first(ls []T) *T {
	if len(ls) == 0 {
		return nil
	}
	return &ls[0]
}

func (slf *Repo[T]) doInsert(tx abs.DbTx, sqlPars ...interface{}) (*srRepoSql[T], error) {
	sql := srRepoSql[T]{}.new(`
INSERT INTO ::tableName (
	::insertColumnNames
) VALUES (
	::insertParamSign
)`)

	sql.pars = sqlPars
	err := sql.transform(slf)
	if err != nil {
		return sql, err
	}

	if tx != nil {
		return sql, slf.DbInstance.TxExecute(tx, sql.query, sql.pars...)
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars...)
}

func (slf *Repo[T]) doInsertRID(tx abs.DbTx, sqlPars ...interface{}) (*srRepoSql[T], *int64, error) {
	sql := srRepoSql[T]{}.new(`
INSERT INTO ::tableName (
	::insertColumnNames
) VALUES (
	::insertParamSign
)
RETURNING id`)

	sql.pars = sqlPars
	err := sql.transform(slf)
	if err != nil {
		return sql, nil, err
	}

	if tx != nil {
		id, err := slf.DbInstance.TxExecuteRID(tx, sql.query, sql.pars...)
		return sql, id, err
	}

	id, err := slf.DbInstance.ExecuteRID(sql.query, sql.pars...)
	return sql, id, err
}

func (slf *Repo[T]) doUpdate(tx abs.DbTx, keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) (*srRepoSql[T], error) {
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

	if tx != nil {
		return sql, slf.DbInstance.TxExecute(tx, sql.query, sql.pars...)
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars...)
}

func (slf *Repo[T]) goGetDatas(tx abs.DbTx, whereQuery string, endQuery string, wherePars ...interface{}) (*srRepoSql[T], []T, error) {
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

	var unsafe *model.DbUnsafeSelectError
	if tx != nil {
		unsafe, err = slf.DbInstance.TxSelect(tx, &models, sql.query, sql.pars)
	} else {
		unsafe, err = slf.DbInstance.Select(&models, sql.query, sql.pars)
	}

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

func (slf *Repo[T]) doSelect(tx abs.DbTx, query string, pars ...interface{}) (*srRepoSql[T], []T, error) {
	var models []T
	sql := srRepoSql[T]{}.new("")

	sql.query = strings.TrimSpace(query)
	sql.pars = pars

	query, pars, err := transformIn(sql.query, sql.pars...)
	if err != nil {
		return sql, models, err
	}

	sql.query = query
	sql.pars = pars

	var unsafe *model.DbUnsafeSelectError
	if tx != nil {
		unsafe, err = slf.DbInstance.TxSelect(tx, &models, sql.query, sql.pars...)
	} else {
		unsafe, err = slf.DbInstance.Select(&models, sql.query, sql.pars...)
	}
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

func (slf *Repo[T]) doCount(tx abs.DbTx, whereQuery, endQuery string, wherePars ...interface{}) (*srRepoSql[T], int, error) {
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
	if tx != nil {
		err = slf.DbInstance.TxGet(tx, &count, sql.query, sql.pars...)
	} else {
		err = slf.DbInstance.Get(&count, sql.query, sql.pars...)
	}

	return sql, count, err
}

func (slf *Repo[T]) doSum(tx abs.DbTx, column, whereQuery, endQuery string, wherePars ...interface{}) (*srRepoSql[T], fc.FCT, error) {
	sql := srRepoSql[T]{}.new(fmt.Sprintf("SELECT COALESCE(SUM(%v), 0) FROM ::tableName", column))

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

	sumValue := fc.New(0)

	err := sql.transform(slf)
	if err != nil {
		return nil, sumValue, err
	}

	if tx != nil {
		err = slf.DbInstance.TxGet(tx, &sumValue, sql.query, sql.pars...)
	} else {
		err = slf.DbInstance.Get(&sumValue, sql.query, sql.pars...)
	}

	return sql, sumValue, err
}

func (slf *Repo[T]) doRawCount(tx abs.DbTx, query string, pars ...interface{}) (*srRepoSql[T], int, error) {
	sql, val, err := slf.doRawFCT(tx, query, pars...)
	return sql, val.Int(), err
}

func (slf *Repo[T]) doRawInt(tx abs.DbTx, query string, pars ...interface{}) (*srRepoSql[T], int, error) {
	sql, val, err := slf.doRawFCT(tx, query, pars...)
	return sql, val.Int(), err
}

func (slf *Repo[T]) doRawInt64(tx abs.DbTx, query string, pars ...interface{}) (*srRepoSql[T], int64, error) {
	sql, val, err := slf.doRawFCT(tx, query, pars...)
	return sql, val.Int64(), err
}

func (slf *Repo[T]) doRawFloat64(tx abs.DbTx, query string, pars ...interface{}) (*srRepoSql[T], float64, error) {
	sql, val, err := slf.doRawFCT(tx, query, pars...)
	return sql, val.Float64(), err
}

func (slf *Repo[T]) doRawFCT(tx abs.DbTx, query string, pars ...interface{}) (*srRepoSql[T], fc.FCT, error) {
	sql := srRepoSql[T]{}.new(``)
	val := fc.New(-1)

	sql.query = strings.TrimSpace(query)
	sql.pars = pars

	query, pars, err := transformIn(sql.query, sql.pars...)
	if err != nil {
		return sql, val, err
	}

	sql.query = query
	sql.pars = pars

	if tx != nil {
		err = slf.DbInstance.TxGet(tx, &val, sql.query, sql.pars...)
	} else {
		err = slf.DbInstance.Get(&val, sql.query, sql.pars...)
	}

	return sql, val, err
}

func (slf *Repo[T]) doDelete(tx abs.DbTx, whereQuery string, wherePars ...interface{}) (*srRepoSql[T], error) {
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

	if tx != nil {
		return sql, slf.DbInstance.TxExecute(tx, sql.query, sql.pars...)
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars...)
}

func (slf *Repo[T]) doExecute(tx abs.DbTx, sqlQuery string, sqlPars ...interface{}) (*srRepoSql[T], error) {
	sql := srRepoSql[T]{}.new(sqlQuery)
	sql.pars = sqlPars

	err := sql.transform(slf)
	if err != nil {
		return nil, err
	}

	if tx != nil {
		return sql, slf.DbInstance.TxExecute(tx, sql.query, sql.pars...)
	}

	return sql, slf.DbInstance.Execute(sql.query, sql.pars...)
}

func sendDbq(logc clog.Instance, sqlQuery string, sqlPars []interface{}, execFunc, execPath string, startAt time.Time, err error) {
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

	clog.SendDbq(0, logc, severity, m, true)
}
