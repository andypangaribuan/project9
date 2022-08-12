/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "google.golang.org/grpc"

type Server interface {
	StartGRPC(port int, autoRecover bool, register func(svc *grpc.Server), stackTraceSkipLevel ...int)
}
