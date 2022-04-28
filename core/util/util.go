/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"github.com/andypangaribuan/project9/abs"
	"github.com/speps/go-hashids/v2"
)

type srEnv struct {
	abs.UtilEnv
}

type srHashId struct {
	abs.UtilHashId
	instance *hashids.HashIDData
}

func Create() (*srEnv, *srHashId) {
	return &srEnv{}, &srHashId{
		instance: getHashIdInstance(),
	}
}
