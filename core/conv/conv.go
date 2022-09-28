/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

type srConv struct {}

type srTime struct{}
type srProto struct{}

func Create() (*srConv, *srTime, *srProto) {
	return &srConv{}, &srTime{}, &srProto{}
}
