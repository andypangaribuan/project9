/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"errors"
	"log"
	"runtime/debug"

	"github.com/shopspring/decimal"
)

// supported operator: +, -, *, /, %
func Cal(val ...interface{}) FCT {
	fcv, err := SCal(val...)
	if err != nil {
		debug.PrintStack()
		objects := []interface{}{err}
		objects = append(objects, val...)
		log.Panicf("error: %+v\nval: %v", objects...)
	}

	return fcv
}

// supported operator: +, -, *, /, %
func SCal(val ...interface{}) (FCT, error) {
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

			if operator == "*" || operator == "/" || operator == "%" {
				vd1 := lsv[i].(decimal.Decimal)
				vd2 := lsv[i+2].(decimal.Decimal)

				switch operator {
				case "*":
					lsv[i] = vd1.Mul(vd2)

				case "/":
					lsv[i] = vd1.Div(vd2)

				case "%":
					lsv[i] = vd1.Mod(vd2)
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
