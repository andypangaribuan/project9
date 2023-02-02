/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"errors"
	"reflect"

	"github.com/shopspring/decimal"
)

func toDecimal(val interface{}) (decimal.Decimal, error) {
	var d decimal.Decimal

	if val == nil {
		return d, errors.New("val cannot nil")
	}

	switch v := val.(type) {
	case string:
		dv, err := decimal.NewFromString(v)
		if err != nil {
			return d, err
		}

		return dv, nil

	case int:
		v64 := int64(v)
		return decimal.NewFromInt(v64), nil

	case int32:
		return decimal.NewFromInt32(v), nil

	case int64:
		return decimal.NewFromInt(v), nil

	case float32:
		return decimal.NewFromFloat32(v), nil

	case float64:
		return decimal.NewFromFloat(v), nil

	case decimal.Decimal:
		return v, nil

	case FCT:
		return v.vd, nil
	}

	if rv := reflect.ValueOf(val); rv.Kind() == reflect.Ptr {
		return toDecimal(rv.Elem().Interface())
	}

	return d, errors.New("unknown type")
}

func isOperation(val interface{}) (string, bool) {
	switch v := val.(type) {
	case string:
		if v == "+" || v == "-" || v == "*" || v == "/" || v == "%" {
			return v, true
		}
	}

	return "", false
}

func removeIndex[T any](ls []T, index int) []T {
	return append(ls[:index], ls[index+1:]...)
}
