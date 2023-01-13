/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func Ptr[T any](val T) *T {
	return &val
}

func PtrNotNil[T any](val *T, fn func() *T) *T {
	if val == nil {
		return nil
	}
	return fn()
}

func PtrNotNilVoid[T any](val *T, fn func()) {
	if val != nil {
		fn()
	}
}

func PtrValue[T any](val *T, defaultValue T) T {
	if val == nil {
		return defaultValue
	}

	return *val
}
