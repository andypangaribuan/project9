/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type srGrpcServerHandler struct {
	stackTraceSkipLevel int
}

func (*srServer) StartGRPC(port int, autoRecover bool, register func(svc *grpc.Server), stackTraceSkipLevel ...int) {
	address := fmt.Sprintf(":%v", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port: 5002\nerror: %v\n", err)
	}

	var server *grpc.Server
	if !autoRecover {
		server = grpc.NewServer()
	} else {
		level := 0
		if len(stackTraceSkipLevel) > 0 && stackTraceSkipLevel[0] > 0 {
			level = stackTraceSkipLevel[0]
		}
		handler := &srGrpcServerHandler{
			stackTraceSkipLevel: level,
		}

		uIntOpt := grpc.UnaryInterceptor(handler.unaryPanicHandler)
		sIntOpt := grpc.StreamInterceptor(handler.streamPanicHandler)
		server = grpc.NewServer(uIntOpt, sIntOpt)
	}

	register(server)

	fmt.Printf("\nGRPC PORT: %v\n", port)
	if err := server.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}

func (slf *srGrpcServerHandler) unaryPanicHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer slf.crashHandler(func(r interface{}) {
		err = slf.toPanicError(r)
	})

	return handler(ctx, req)
}

func (slf *srGrpcServerHandler) streamPanicHandler(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer slf.crashHandler(func(r interface{}) {
		err = slf.toPanicError(r)
	})

	return handler(srv, stream)
}

func (slf *srGrpcServerHandler) toPanicError(r interface{}) error {
	return status.Errorf(codes.Internal, "panic: %v", r)
}

func (slf *srGrpcServerHandler) crashHandler(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
		slf.printPanic(r)
	}
}

func (slf *srGrpcServerHandler) printPanic(r interface{}) {
	var callers []string

	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		callers = append(callers, fmt.Sprintf("%d: %v:%v\n", i, file, line))
	}

	if len(callers) == 0 {
		return
	}

	fmt.Printf("## recovered from panic\n")
	fmt.Printf("## detail\n[[%#v]]\n", r)
	fmt.Printf("## stacktrace:\n")

	startIndex := 0
	if slf.stackTraceSkipLevel > 0 {
		if slf.stackTraceSkipLevel < len(callers) {
			startIndex = slf.stackTraceSkipLevel
		}
	}

	for i := startIndex; len(callers) > i; i++ {
		fmt.Printf(" %v", callers[i])
	}

	fmt.Println()
}
