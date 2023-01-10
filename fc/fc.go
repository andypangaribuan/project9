/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"errors"
	"log"
	"strconv"

	"github.com/shopspring/decimal"
)

func FCal(val ...interface{}) FCT {
	fcv, err := FSCal(val...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	return fcv
}

func FSCal(val ...interface{}) (FCT, error) {
	fv := FCT{
		V1: "0",
		V2: "0",
	}

	length := len(val)

	if length%2 == 0 || length == 0 {
		return fv, errors.New("wrong implementation")
	}

	if length == 1 {
		vd, err := toDecimal(val[0])
		if err != nil {
			return fv, err
		}

		fv.set(vd)
		return fv, nil
	}

	lsv := make([]interface{}, 0)

	for i := 0; i < length; i++ {
		if i%2 == 0 {
			vd, err := toDecimal(val[i])
			if err != nil {
				return fv, err
			}

			lsv = append(lsv, vd)
		} else {
			operation, ok := isOperation(val[i])
			if !ok {
				return fv, errors.New("wrong implementation: invalid operation")
			}

			lsv = append(lsv, operation)
		}
	}

	for i := 0; i < len(lsv); i++ {
		if i%2 == 0 && i < len(lsv)-1 {
			operator := lsv[i+1].(string)

			if operator == "*" || operator == "/" {
				vd1 := lsv[i].(decimal.Decimal)
				vd2 := lsv[i+2].(decimal.Decimal)

				switch operator {
				case "*":
					lsv[i] = vd1.Mul(vd2)

				case "/":
					lsv[i] = vd1.Div(vd2)
				}

				lsv = removeIndex(lsv, i+2)
				lsv = removeIndex(lsv, i+1)
				i--
			}
		}
	}

	for i := 0; i < len(lsv); i++ {
		if i%2 == 0 && i < len(lsv)-1 {
			operator := lsv[i+1].(string)

			if operator == "+" || operator == "-" {
				vd1 := lsv[i].(decimal.Decimal)
				vd2 := lsv[i+2].(decimal.Decimal)

				switch operator {
				case "+":
					lsv[i] = vd1.Add(vd2)

				case "-":
					lsv[i] = vd1.Sub(vd2)
				}

				lsv = removeIndex(lsv, i+2)
				lsv = removeIndex(lsv, i+1)
				i--
			}
		}
	}

	if len(lsv) != 1 {
		return fv, errors.New("something went wrong")
	}

	fv.set(lsv[0].(decimal.Decimal))
	return fv, nil
}

func toDecimal(val interface{}) (decimal.Decimal, error) {
	var d decimal.Decimal

	switch v := val.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return d, err
		}

		d = decimal.NewFromFloat(f)

		return d, nil

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

	return d, errors.New("unknown type")
}

func isOperation(val interface{}) (string, bool) {
	switch v := val.(type) {
	case string:
		if v == "+" || v == "-" || v == "*" || v == "/" {
			return v, true
		}
	}

	return "", false
}

func removeIndex[T any](ls []T, index int) []T {
	return append(ls[:index], ls[index+1:]...)
}
