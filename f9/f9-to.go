/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"fmt"
	"strconv"
	"strings"
)

 func ToCurrencyFormat[T string | int | int32 | int64 | float32 | float64](val T, decimalPoint int, thousandSeparator, decimalSeparator string) string {
	currency := ""

	convert := func(val float64) {
		format := "%." + fmt.Sprintf("%v", decimalPoint) + "f"
		v := fmt.Sprintf(format, val)
		ls := strings.Split(v, ".")

		l0 := ""
		count := 0
		ls0Length := len(ls[0])
		for i := ls0Length - 1; i >= 0; i-- {
			l0 = ls[0][i:i+1] + l0
			count++
			if i > 0 && count == 3 {
				count = 0
				l0 = thousandSeparator + l0
			}
		}

		currency = l0
		if len(ls) > 1 {
			currency += decimalSeparator + ls[1]
		}
	}

	var iv interface{}
	iv = val

	switch iv := iv.(type) {
	case int:
		convert(float64(iv))
	case int32:
		convert(float64(iv))
	case int64:
		convert(float64(iv))
	case float32:
		convert(float64(iv))
	case float64:
		convert(iv)
	case string:
		v, err := strconv.Atoi(iv)
		if err == nil {
			convert(float64(v))
		} else {
			v, err := strconv.ParseFloat(iv, 64)
			if err == nil {
				convert(v)
			}
		}
	}

	return currency
}