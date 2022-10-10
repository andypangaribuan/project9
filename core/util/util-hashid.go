/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/p9"
	"github.com/speps/go-hashids/v2"
)

func getHashIdInstance(salt string, length int) *hashids.HashIDData {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	return hd
}

func (slf *srHashId) Reload() {
	slf.instance = getHashIdInstance(p9.Conf.HashIdSalt, p9.Conf.HashIdLength)
}

func (slf *srHashId) Encode(numbers ...int) string {
	hid, _ := hashids.NewWithData(slf.instance)
	hashed, _ := hid.Encode(numbers)
	return hashed
}

func (slf *srHashId) Decode(hashed string) ([]int, error) {
	hid, _ := hashids.NewWithData(slf.instance)
	return hid.DecodeWithError(hashed)
}

func (slf *srHashId) Encode64(numbers ...int64) string {
	hid, _ := hashids.NewWithData(slf.instance)
	hashed, _ := hid.EncodeInt64(numbers)
	return hashed
}

func (slf *srHashId) DecodeInt64(hashed string) ([]int64, error) {
	hid, _ := hashids.NewWithData(slf.instance)
	return hid.DecodeInt64WithError(hashed)
}

func (slf *srHashId) Add(key, salt string, length int) {
	sr := &srHashId{
		instance: getHashIdInstance(salt, length),
		slfMap:   make(map[string]*srHashId, 0),
	}
	slf.slfMap[key] = sr
}

func (slf *srHashId) Get(key string) abs.UtilHashId {
	return slf.slfMap[key]
}
