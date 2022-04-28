/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

import "github.com/andypangaribuan/project9/abs"

type srTime struct {
	abs.ConvTime
}

func Create() *srTime {
	return &srTime{}
}
