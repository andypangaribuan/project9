/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package test

import (
	"fmt"
	"testing"

	"github.com/andypangaribuan/project9/fc"
	"github.com/shopspring/decimal"
)

func TestPerimeter(t *testing.T) {
	fcv := fc.FCal(0.1, "*", 0.2, "/", 0.2) //, "+", 0.3, "+", 0.2, "*", 0.3)
	fcv = fc.FCal(fcv, "+", 0.3)
	fcv = fc.FCal(fcv, "+", 0.2, "*", 0.3)
	d1 := decimal.NewFromFloat(0.1)
	d2 := decimal.NewFromFloat(0.2)
	d3 := decimal.NewFromFloat(0.2)
	d4 := decimal.NewFromFloat(0.3)
	d5 := decimal.NewFromFloat(0.2)
	d6 := decimal.NewFromFloat(0.3)
	d := d1.Mul(d2).Div(d3).Add(d4).Add(d5.Mul(d6))

	fmt.Printf("v: %.20f, e: %.20f, r: %.20f\n", fcv.Float64(), d.InexactFloat64(), 0.1*0.2/0.2+0.3+0.2*0.3)

	v1 := 0.3
	v2 := 0.1
	v3 := 10000.0
	d = deci(v1).Div(deci(v2)).Mul(deci(v3))
	v := v1 / v2 * v3
	f := fc.FCal(v1, "/", v2, "*", v3)
	fmt.Printf("d: %.20f, f: %.20f, v: %.20f, r: %.20f, r1: %.20f\n", d.InexactFloat64(), f.Float64(), v, 0.3/0.1*10000, 3.0)
}

func deci(v interface{}) decimal.Decimal {
	switch val := v.(type) {
	case int:
		return decimal.NewFromInt(int64(val))
	case int32:
		return decimal.NewFromInt32(val)
	case int64:
		return decimal.NewFromInt(val)
	case float32:
		return decimal.NewFromFloat32(val)
	case float64:
		return decimal.NewFromFloat(val)
	}

	return decimal.NewFromInt(0)
}
