/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func Ptr[T any](val T) *T {
	return &val
}
