/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import (
	"github.com/andypangaribuan/project9/abs"
)

type srConv struct {
	Time  *srConvTime
	Proto *srConvProto
}

type srConvTime struct {
	absConvTime
}

type absConvTime interface {
	abs.ConvTime
}

type srConvProto struct {
	absConvProto
}

type absConvProto interface {
	abs.ConvProto
}
