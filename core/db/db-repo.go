/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
	"github.com/andypangaribuan/project9/p9"
	"github.com/pkg/errors"
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
	slf.ColumnNames = strings.TrimSpace(names)
}

func (slf *Repo[T]) SetInsertColumnNames(names string) {
	slf.InsertColumnNames = strings.TrimSpace(names)
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

func (slf *Repo[T]) Save(sqlPars ...interface{}) error {
	sql := srRepoSql[T]{}.new(`
INSERT INTO ::tableName (
	::insertColumnNames
) VALUES (
	::insertParamSign
)`)
	sql.transform(slf)

	return slf.DbInstance.Execute(sql.query, sqlPars)
}

func (slf *Repo[T]) Update(keyVals map[string]interface{}, whereQuery string, wherePars ...interface{}) error {
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

	sql.transform(slf)

	return slf.DbInstance.Execute(sql.query, sql.pars)
}

func (slf *Repo[T]) GetDatas(whereQuery string, endQuery string, wherePars ...interface{}) ([]T, error) {
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
	sql.transform(slf)
	sql.pars = wherePars

	unsafe, err := slf.DbInstance.Select(&models, sql.query, sql.pars)
	slf.OnUnsafe(unsafe)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return models, nil
}

func (slf *Repo[T]) GetData(whereQuery string, endQuery string, wherePars ...interface{}) (*T, error) {
	models, err := slf.GetDatas(whereQuery, endQuery, wherePars...)
	return slf.first(models), err
}
