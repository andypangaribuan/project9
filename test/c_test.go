/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package test

import (
	"fmt"
	"testing"
)

func TestC(t *testing.T) {
	size := 10

	c1 := 10 % size
	fmt.Printf("c1: %v\n", c1)

	c2 := 11 % size
	fmt.Printf("c1: %v\n", c2)

	c3 := 19 % size
	fmt.Printf("c1: %v\n", c3)

	f1 := 10 / size
	fmt.Printf("f1: %v\n", f1)

	f2 := 11 / size
	fmt.Printf("f2: %v\n", f2)

	f3 := 19 / float64(size)
	f3s := fmt.Sprintf("%.1f", f3)
	fmt.Printf("f3: %v\n", f3s)

	f4 := 20 / float64(size)
	f4s := fmt.Sprintf("%.1f", f4)
	fmt.Printf("f3: %v\n", f4s)

	f5 := 0 / float64(size)
	f5s := fmt.Sprintf("%.1f", f5)
	fmt.Printf("f3: %v\n", f5s)
}