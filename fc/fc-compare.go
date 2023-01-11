/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"errors"
	"log"

	"github.com/andypangaribuan/project9/f9"
)

func Compare(v1 interface{}, operation string, v2 interface{}) bool {
	v, err := SCompare(v1, operation, v2)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	return v
}

func SCompare(v1 interface{}, operation string, v2 interface{}) (bool, error) {
	if !f9.IfHaveIn(operation, "==", "<", "<=", ">=", ">") {
		return false, errors.New("invalid operation")
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
