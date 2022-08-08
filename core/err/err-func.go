/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package err

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/pkg/errors"
)

type stack []uintptr

func (*srErr) WithStack(err error, skip ...int) error {
	if err == nil {
		return err
	}

	skipLevel := 3
	if len(skip) > 0 {
		skipLevel += skip[0]
	}

	val := reflect.ValueOf(err)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	vType := val.Type()
	for i := 0; i < vType.NumField(); i++ {
		sf := vType.Field(i)
		if sf.PkgPath == "github.com/pkg/errors" {
			return err
		}
	}

	err = errors.WithStack(err)
	pointerVal := reflect.ValueOf(err)
	val = reflect.Indirect(pointerVal)
	vType = val.Type()

	for i := 0; i < vType.NumField(); i++ {
		sf := vType.Field(i)
		fv := val.Field(i)
		if sf.Name == "stack" {
			ptrToY := unsafe.Pointer(fv.UnsafeAddr())
			realPtrToY := (**stack)(ptrToY)
			*realPtrToY = callers(skipLevel)
			break
		}
	}

	return err
}

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
