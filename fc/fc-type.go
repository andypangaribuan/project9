/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/shopspring/decimal"
)

type FCT struct {
	vd decimal.Decimal
	V1 string
	V2 string
}

func New(val interface{}) FCT {
	fv, err := SNew(val)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	return fv
}

func SNew(val interface{}) (FCT, error) {
	var fv FCT

	dv, err := toDecimal(val)
	if err != nil {
		return fv, nil
	}

	fv.set(dv)
	return fv, nil
}

func (slf *FCT) set(vd decimal.Decimal) {
	slf.vd = vd
	slf.V1 = fmt.Sprintf("%.20f", slf.vd.InexactFloat64())
	slf.V2 = printer.Sprintf("%.20f", slf.vd.InexactFloat64())
}

func (slf FCT) Float64() float64 {
	return slf.vd.InexactFloat64()
}

func (slf FCT) Round(places int) FCT {
	return New(slf.vd.Round(int32(places)))
}

func (slf FCT) Int() int {
	return int(slf.vd.IntPart())
}

func (slf FCT) Int64() int64 {
	return slf.vd.IntPart()
}

func (slf FCT) Decimal() decimal.Decimal {
	return slf.vd
}

func (slf FCT) Floor(places ...int) FCT {
	if len(places) == 0 {
		return New(slf.vd.Floor())
	}

	if places[0] < 1 {
		return New(slf.vd.Floor())
	}

	dc1 := 1
	for i := 0; i < places[0]; i++ {
		v, err := strconv.Atoi(fmt.Sprintf("%v0", dc1))
		if err != nil {
			log.Fatalf("converter error: %+v\n", err)
		}

		dc1 = v
	}

	dc2 := decimal.NewFromInt(int64(dc1))
	return New(slf.vd.Mul(dc2).Floor().Div(dc2))
}
