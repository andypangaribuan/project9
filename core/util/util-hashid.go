/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"github.com/andypangaribuan/project9/conf"
	"github.com/speps/go-hashids/v2"
)

func getHashIdInstance() *hashids.HashIDData {
	hd := hashids.NewData()
	hd.Salt = conf.HashIdSalt
	hd.MinLength = conf.HashIdLength
	return hd
}

func (slf *srHashId) Reload() {
	slf.instance = getHashIdInstance()
}

func (slf *srHashId) Encode(numbers []int) string {
	hid, _ := hashids.NewWithData(slf.instance)
	hashed, _ := hid.Encode(numbers)
	return hashed
}

func (slf *srHashId) EncodeId(number int) string {
	return slf.Encode([]int{number})
}
