/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type srRepoSql[T any] struct {
	query string
	pars  []interface{}
}

func (slf srRepoSql[T]) new(query string) *srRepoSql[T] {
	slf.query = query
	slf.pars = make([]interface{}, 0)
	return &slf
}

func (slf *srRepoSql[T]) transform(base *Repo[T]) error {
	insertColumnNames := ""

	if strings.Contains(slf.query, "::tableName") {
		slf.query = strings.ReplaceAll(slf.query, "::tableName", base.TableName)
	}

	if strings.Contains(slf.query, "::insertColumnNames") {
		insertColumnNames = base.InsertColumnNames
		if insertColumnNames == "" {
			insertColumnNames = base.ColumnNames
		}
		insertColumnNames = strings.TrimSpace(insertColumnNames)
		slf.query = strings.ReplaceAll(slf.query, "::insertColumnNames", insertColumnNames)
	}

	if strings.Contains(slf.query, "::insertParamSign") {
		insertParamSign := base.InsertParamSign
		if insertParamSign == "" {
			insertParamSign = base.ParamSigns
		}
		if insertParamSign == "" && insertColumnNames != "" {
			icn := strings.ReplaceAll(insertColumnNames, " ", "")
			arr := strings.Split(icn, ",")
			for i := range arr {
				if i > 0 {
					insertParamSign += ", "
				}
				insertParamSign += "?"
			}
		}
		insertParamSign = strings.TrimSpace(insertParamSign)
		slf.query = strings.ReplaceAll(slf.query, "::insertParamSign", insertParamSign)
	}

	slf.query = strings.TrimSpace(slf.query)

	if strings.Contains(strings.ToUpper(slf.query), " IN ") {
		query, pars, err := transformIn(slf.query, slf.pars...)
		if err != nil {
			return err
		}

		slf.query = query
		slf.pars = pars
	}

	return nil
}

func transformIn(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}
