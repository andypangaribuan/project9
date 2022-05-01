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
	OnUnsafe          func(unsafe *model.DbUnsafeSelectError)
}

func (slf *Repo[T]) onUnsafe(unsafe *model.DbUnsafeSelectError) {
	if unsafe != nil && slf.PrintUnsafeErr {
		fmt.Printf("[%v] db.unsafe.select.error:\nerror: %v\nsql-query: %v\nsql-pars: %v\ntrace: %v\n",
			p9.Conv.Time.ToStrFull(time.Now()), unsafe.LogMessage, unsafe.SqlQuery, unsafe.SqlPars, unsafe.LogTrace)
	}
}

func (slf *Repo[T]) first(ls []T) *T {
	if len(ls) == 0 {
		return nil
	}
	return &ls[0]
}

func (slf *Repo[T]) GenerateParamSigns(columnNames string) string {
	cn := strings.TrimSpace(columnNames)
	cn = strings.ReplaceAll(cn, " ", "")
	ls := strings.Split(cn, ",")
	paramSign := ""
	for i := range ls {
		if i > 0 {
			paramSign += ","
		}
		paramSign += " ?"
	}
	return paramSign
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

func (slf *Repo[T]) GetDatas(whereQuery string, wherePars ...interface{}) ([]T, error) {
	var models []T
	sql := srRepoSql[T]{}.new(`SELECT * FROM ::tableName`)

	whereQuery = strings.TrimSpace(whereQuery)
	wq := strings.ToLower(whereQuery)
	if (len(wq) <= 6) || (len(wq) > 6 && wq[:6] != "where ") {
		whereQuery = fmt.Sprintf("WHERE %v", whereQuery)
	}

	sql.query += fmt.Sprintf(" %v", whereQuery)
	sql.transform(slf)

	unsafe, err := slf.DbInstance.Select(&models, sql.query, wherePars)
	slf.onUnsafe(unsafe)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return models, nil
}

func (slf *Repo[T]) GetData(whereQuery string, wherePars ...interface{}) (*T, error) {
	models, err := slf.GetDatas(whereQuery, wherePars...)
	return slf.first(models), err
}
