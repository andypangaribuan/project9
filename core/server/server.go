/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

type srServer struct{}

func Create() *srServer {
	return &srServer{}
}
