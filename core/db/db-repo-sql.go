/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import "strings"

type srRepoSql[T any] struct {
	query string
	pars  []interface{}
}

func (slf srRepoSql[T]) new(query string) *srRepoSql[T] {
	slf.query = query
	slf.pars = make([]interface{}, 0)
	return &slf
}

func (slf *srRepoSql[T]) transform(base *Repo[T]) {
	if strings.Contains(slf.query, "::tableName") {
		slf.query = strings.ReplaceAll(slf.query, "::tableName", base.TableName)
	}

	if strings.Contains(slf.query, "::insertColumnNames") {
		insertColumnNames := base.InsertColumnNames
		if insertColumnNames == "" {
			insertColumnNames = base.ColumnNames
		}
		slf.query = strings.ReplaceAll(slf.query, "::insertColumnNames", insertColumnNames)
	}

	if strings.Contains(slf.query, "::insertParamSign") {
		insertParamSign := base.InsertParamSign
		if insertParamSign == "" {
			insertParamSign = base.ParamSigns
		}
		slf.query = strings.ReplaceAll(slf.query, "::insertParamSign", insertParamSign)
	}

	slf.query = strings.TrimSpace(slf.query)
}
