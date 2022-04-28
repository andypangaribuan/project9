/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"github.com/andypangaribuan/project9/p9"
	"github.com/matoous/go-nanoid/v2"
)

func (*srUtil) GetNanoID(length ...int) (string, error) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	size := p9.Conf.NanoIdLength
	if len(length) > 0 {
		size = length[0]
	}

	return gonanoid.Generate(alphabet, size)
}
