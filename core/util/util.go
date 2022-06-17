/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"github.com/andypangaribuan/project9/p9"
	"github.com/speps/go-hashids/v2"
)

type srUtil struct{}

type srEnv struct{}

type srEnvBase64 struct {
	key  string
	data []byte
}

type srHashId struct {
	instance *hashids.HashIDData
	slfMap   map[string]*srHashId
}

func Create() (*srUtil, *srEnv, *srHashId) {
	return &srUtil{}, &srEnv{}, &srHashId{
		instance: getHashIdInstance(p9.Conf.HashIdSalt, p9.Conf.HashIdLength),
		slfMap:   make(map[string]*srHashId, 0),
	}
}
