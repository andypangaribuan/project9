/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"reflect"
	"strings"

	"github.com/andypangaribuan/project9/abs"
)

type srUpdateHelper struct {
	columnSet   []string
	columnWhere []string
	sqlWhere    string
	parsSet     []interface{}
	parsWhere   []interface{}
}

func NewUpdateHelper() abs.DbUpdateHelper {
	v := &srUpdateHelper{
		columnSet:   make([]string, 0),
		columnWhere: make([]string, 0),
		parsSet:     make([]interface{}, 0),
		parsWhere:   make([]interface{}, 0),
	}

	return v
}

func (slf *srUpdateHelper) SetAdd(condition string, param interface{}) {
	slf.columnSet = append(slf.columnSet, condition)
	slf.parsSet = append(slf.parsSet, param)
}

func (slf *srUpdateHelper) SetAddIfNoNil(condition string, param interface{}) {
	if param != nil {
		slf.columnSet = append(slf.columnSet, condition)
		slf.parsSet = append(slf.parsSet, param)
	}
}

func (slf *srUpdateHelper) SetAddIfNotNilOrEmpty(condition string, param interface{}) {
	if param != nil {
		isAdd := false

		if rv := reflect.ValueOf(param); rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return
			}

			param = rv.Elem().Interface()
		}

		switch val := param.(type) {
		case string:
			if len(strings.TrimSpace(val)) > 0 {
				isAdd = true
			}

		default:
			if param != nil {
				isAdd = true
			}
		}

		if isAdd {
			slf.columnSet = append(slf.columnSet, condition)
			slf.parsSet = append(slf.parsSet, param)
		}
	}
}

func (slf *srUpdateHelper) Where(condition string, pars ...interface{}) {
	slf.sqlWhere = condition
	slf.parsWhere = pars
}

func (slf *srUpdateHelper) MapSetAdd(kv map[string]interface{}) {
	for condition, param := range kv {
		slf.SetAdd(condition, param)
	}
}

func (slf *srUpdateHelper) MapSetAddIfNotNilOrEmpty(kv map[string]interface{}) {
	for condition, param := range kv {
		slf.SetAddIfNotNilOrEmpty(condition, param)
	}
}

func (slf *srUpdateHelper) Get() (condition string, pars []interface{}) {
	sqlSet := ""
	for i, condition := range slf.columnSet {
		if i == 0 {
			sqlSet = "SET "
		} else {
			sqlSet += ", "
		}

		sqlSet += condition
	}

	if sqlSet == "" {
		return
	}

	sqlWhere := ""
	if slf.sqlWhere != "" {
		sqlWhere = "WHERE " + slf.sqlWhere
	}

	condition = sqlSet + " " + sqlWhere
	pars = append(pars, slf.parsSet...)
	pars = append(pars, slf.parsWhere...)
	return
}
