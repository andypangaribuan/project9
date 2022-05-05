/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	} else {
		return b
	}
}

func TernaryFn[T any](condition bool, a, b func() T) T {
	if condition {
		return a()
	} else {
		return b()
	}
}

func TernaryFnA[T any](condition bool, a func() T, b T) T {
	if condition {
		return a()
	} else {
		return b
	}
}

func TernaryFnB[T any](condition bool, a T, b func() T) T {
	if condition {
		return a
	} else {
		return b()
	}
}
