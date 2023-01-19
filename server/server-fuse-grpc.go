/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/server/proto/gen/grf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (slf *srFuseGrpc) Restful(c context.Context, req *grf.Request) (*grf.Response, error) {
	var (
		err error
		ctx = &srFuseContext{
			path: req.Action,
			grpcCtx: &srFuseGrpcContext{
				ctx:     c,
				request: req,
				header:  req.Header,
				payload: req.Payload,
				params:  req.Params,
				queries: req.Queries,
			},
		}
	)

	ctx.reqCtx = &srFuseContextRequest{
		fuseCtx: ctx,
	}

	handlers, ok := slf.routes[req.Action]
	if !ok {
		err = fmt.Errorf(`action not found: "%v"`, req.Action)
		return ctx.grpcSend(ctx.R500InternalServerError(nil, err))
	}

	if len(handlers) >= 2 {
		err = handlers[0](ctx)
		if !ctx.isAuthSet {
			return ctx.grpcSend(err)
		}

		err = handlers[1](ctx)
	} else {
		err = handlers[0](ctx)
	}

	return ctx.grpcSend(err)
}

func ConnectFuseGRPC(address string) (ClientFuseGRPC, error) {
	fuse := &srConnectFuseGRPC{
		address: address,
	}

	err := fuse.buildConnection()
	if err != nil {
		return nil, err
	}

	fuse.buildClient()

	return &srClientFuseGRPC{fuseGrpc: fuse}, nil
}

func (slf *srConnectFuseGRPC) buildConnection() error {
	if slf.conn != nil {
		return nil
	}

	slf.connMutex.Lock()
	defer slf.connMutex.Unlock()
	if slf.conn != nil {
		return nil
	}

	conn, err := grpc.Dial(slf.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	} else {
		// make sure we connected to the service
		timeout := 5 * time.Second
		conn, err := net.DialTimeout("tcp", slf.address, timeout)
		if err != nil {
			return err
		}
		_ = conn.Close()
	}

	slf.conn = conn
	return nil
}

func (slf *srConnectFuseGRPC) buildClient() {
	if slf.client != nil {
		return
	}

	slf.clientMutex.Lock()
	defer slf.clientMutex.Unlock()

	slf.client = grf.NewRestfulServiceClient(slf.conn)
}

func (slf *srClientFuseGRPC) Restful(path string, header map[string]string, payload map[string]interface{}, params map[string]string) (*grf.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	request := &grf.Request{
		Action:  path,
		Header:  header,
		Payload: f9.GetMapStructPbValue(payload),
		Params:  params,
	}

	res, err := slf.fuseGrpc.client.Restful(ctx, request)
	return res, err
}
