/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func Ls[T any](vals ...T) []T {
	ls := make([]T, 0)

	if len(vals) > 0 {
		ls = append(ls, vals...)
	}

	return ls
}
