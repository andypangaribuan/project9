/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import (
	"github.com/andypangaribuan/project9/abs"
)

type srConv struct {
	Time *srConvTime
}

type srConvTime struct {
	abs.ConvTime
}
