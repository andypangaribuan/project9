/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import "fmt"

func FnIfPtrNotNil[T any](val *T, fn ...func(v T)) {
	if val != nil {
		for _, f := range fn {
			f(*val)
		}
	}
}

func FnGO(funcs ...func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("FnGO recover, error: %+v", r)
		}
	}()

	for _, f := range funcs {
		err = f()
		if err != nil {
			break
		}
	}

	return
}
