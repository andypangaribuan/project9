/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/andypangaribuan/project9/f9"
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

	if val == nil {
		return fv, errors.New("val cannot nil")
	}

	dv, err := toDecimal(val)
	if err != nil {
		return fv, nil
	}

	fv.set(dv)
	return fv, nil
}

func (slf *FCT) set(vd decimal.Decimal) {
	exp := int(vd.Exponent())
	if exp < 0 {
		exp *= -1
	}

	if exp < 1 {
		exp = 1
	}
	
	format := "%." + strconv.Itoa(exp) + "f"

	slf.vd = vd
	slf.V1 = fmt.Sprintf(format, slf.vd.InexactFloat64())
	slf.V2 = printer.Sprintf(format, slf.vd.InexactFloat64())
}

func (slf FCT) Float64() float64 {
	return slf.vd.InexactFloat64()
}

func (slf *FCT) PtrFloat64(defaultValue ...float64) *float64 {
	if slf != nil {
		return f9.Ptr(slf.vd.InexactFloat64())
	}

	if len(defaultValue) > 0 {
		return &defaultValue[0]
	}

	return nil
}

func (slf FCT) Round(places int) FCT {
	return New(slf.vd.Round(int32(places)))
}

func (slf *FCT) PtrRound(places int, defaultValue ...FCT) *FCT {
	if slf != nil {
		return f9.Ptr(New(slf.vd.Round(int32(places))))
	}

	if len(defaultValue) > 0 {
		return &defaultValue[0]
	}

	return nil
}

func (slf FCT) Int() int {
	return int(slf.vd.IntPart())
}

func (slf *FCT) PtrInt(defaultValue ...int) *int {
	if slf != nil {
		return f9.Ptr(slf.Int())
	}

	if len(defaultValue) > 0 {
		return &defaultValue[0]
	}

	return nil
}

func (slf FCT) Int64() int64 {
	return slf.vd.IntPart()
}

func (slf *FCT) PtrInt64(defaultValue ...int64) *int64 {
	if slf != nil {
		return f9.Ptr(slf.vd.IntPart())
	}

	if len(defaultValue) > 0 {
		return &defaultValue[0]
	}

	return nil
}

func (slf FCT) Decimal() decimal.Decimal {
	return slf.vd
}

func (slf *FCT) PtrDecimal() *decimal.Decimal {
	if slf != nil {
		return f9.Ptr(slf.vd)
	}

	return nil
}

func (slf FCT) Floor(places ...int) FCT {
	if len(places) == 0 {
		return New(slf.vd.Floor())
	}

	if places[0] < 1 {
		return New(slf.vd.Floor())
	}

	exp := slf.vd.Exponent()
	if exp < 0 {
		exp *= -1
		if exp > int32(places[0]) {
			sub := int(exp) - places[0]
			div := "1"
			thousandDivDecimal := big.NewInt(1)

			for i := 0; i < sub; i++ {
				div = fmt.Sprintf("%v0", div)
				v, ok := new(big.Int).SetString(div, 10)
				if !ok {
					log.Fatalf("error when converting to big.int, value: %v\n", div)
				}

				thousandDivDecimal = v
			}

			currentValue := slf.vd.Coefficient()
			newValue := new(big.Int).Div(currentValue, thousandDivDecimal)

			return New(decimal.NewFromBigInt(newValue, int32(places[0] * -1)))
		}

	}

	return slf
}

// places parameter using int/int32/int64 type
// defaultValue using fc.FCT type
func (slf *FCT) PtrFloor(opt ...interface{}) *FCT {
	var (
		places       = make([]int, 0)
		defaultValue *FCT
	)

	for _, o := range opt {
		switch v := o.(type) {
		case int:
			if len(places) == 0 {
				places = append(places, v)
			}

		case int32:
			if len(places) == 0 {
				places = append(places, int(v))
			}

		case int64:
			if len(places) == 0 {
				places = append(places, int(v))
			}

		case FCT:
			if defaultValue == nil {
				defaultValue = f9.Ptr(v)
			}
		}
	}

	if slf != nil {
		return f9.Ptr(slf.Floor(places...))
	}

	if defaultValue != nil {
		return f9.Ptr(defaultValue.Floor(places...))
	}

	return nil
}


func (slf *FCT) ToString() string {
	exp := int(slf.vd.Exponent())
	if exp < 0 {
		exp *= -1
	}

	if exp == 0 {
		exp = 1
	}

	format := "%." + strconv.Itoa(exp) + "f"
	return fmt.Sprintf(format, slf.vd.InexactFloat64())
}