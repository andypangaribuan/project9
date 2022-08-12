/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"errors"
	"fmt"
	"reflect"
)

func PanicCatcher(fn func()) (err error) {
	defer func() {
		pv := recover()
		if pv != nil {
			if v, ok := pv.(error); ok {
				err = v
			} else {
				rv := reflect.ValueOf(pv)
				if rv.Kind() == reflect.Ptr {
					pv = rv.Elem()
				}

				msg := fmt.Sprintf("%+v", pv)
				err = errors.New(msg)
			}
		}
	}()
	fn()
	return
}
