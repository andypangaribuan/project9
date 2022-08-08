/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func GetDefault[T any](ls []T, defVal T) T {
	if len(ls) > 0 {
		return ls[0]
	}
	return defVal
}
