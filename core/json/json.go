/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package json

import "github.com/andypangaribuan/project9/abs"

type srJson struct {
	abs.Json
}

func Create() *srJson {
	return &srJson{}
}
