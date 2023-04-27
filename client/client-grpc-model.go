/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package client

import "google.golang.org/grpc"

type GrpcClient[T any] struct {
	address       string
	usingClientLB bool
	usingTLS      bool
	conn          *grpc.ClientConn
	newClientFunc func(*grpc.ClientConn) T
	client        T
}
