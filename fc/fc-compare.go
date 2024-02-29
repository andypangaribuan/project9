/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"errors"
	"log"
	"reflect"
	"runtime/debug"

	"github.com/andypangaribuan/project9/f9"
)

func Compare(v1 interface{}, operation string, v2 interface{}) bool {
	v, err := SCompare(v1, operation, v2)
	if err != nil {
		debug.PrintStack()
		log.Panicf("error: %+v\nv1: %v, op: %v, v2: %v\n", err, v1, operation, v2)
	}

	return v
}

func SCompare(v1 interface{}, operation string, v2 interface{}) (bool, error) {
	if !f9.IfHaveIn(operation, "==", "!=", "<", "<=", ">=", ">") {
		return false, errors.New("fc.SCompare: invalid operation")
	}

	if v1 == nil {
		return false, errors.New("v1 cannot nil")
	}

	if v2 == nil {
		return false, errors.New("v2 cannot nil")
	}

	if rv := reflect.ValueOf(v1); rv.Kind() == reflect.Ptr {
		v1 = rv.Elem().Interface()
	}

	if rv := reflect.ValueOf(v2); rv.Kind() == reflect.Ptr {
		v2 = rv.Elem().Interface()
	}

	var (
		fv1 FCT
		fv2 FCT
	)

	switch v := v1.(type) {
	case FCT:
		fv1 = v

	default:
		fv1 = New(v)
	}

	switch v := v2.(type) {
	case FCT:
		fv2 = v
	default:
		fv2 = New(v)
	}

	switch operation {
	case "==":
		return fv1.vd.Equal(fv2.vd), nil

	case "!=":
		return !fv1.vd.Equal(fv2.vd), nil

	case "<":
		return fv1.vd.LessThan(fv2.vd), nil

	case "<=":
		return fv1.vd.LessThanOrEqual(fv2.vd), nil

	case ">=":
		return fv1.vd.GreaterThanOrEqual(fv2.vd), nil

	case ">":
		return fv1.vd.GreaterThan(fv2.vd), nil
	}

	return false, errors.New("unknown error")
}
