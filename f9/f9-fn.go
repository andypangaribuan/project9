/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func FnIfPtrNotNil[T any](val *T, fn ...func(v T)) {
	if val != nil {
		for _, f := range fn {
			f(*val)
		}
	}
}
