/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type FCT struct {
	vd decimal.Decimal
	V1 string
	V2 string
}

func (slf *FCT) set(vd decimal.Decimal) {
	slf.vd = vd
	slf.V1 = fmt.Sprintf("%.20f", slf.vd.InexactFloat64())
	slf.V2 = printer.Sprintf("%.20f", slf.vd.InexactFloat64())
}

func (slf FCT) Float64() float64 {
	return slf.vd.InexactFloat64()
}
