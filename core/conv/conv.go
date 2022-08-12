/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

type srTime struct{}
type srProto struct{}

func Create() (*srTime, *srProto) {
	return &srTime{}, &srProto{}
}
