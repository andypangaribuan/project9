/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"context"
	"mime/multipart"
	"sync"

	"github.com/andypangaribuan/project9/server/proto/gen/grf"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type FuseContext interface {
	Params(key string) string
	Parser(header, body interface{}) (bool, error)

	AuthX() interface{}
	AuthY() interface{}
	AuthZ() interface{}
	SetAuth(authX, authY, authZ interface{})

	R200OK(data interface{}, opt ...FuseOpt) error
	R400BadRequest(message string, opt ...FuseOpt) error
	R401Unauthorized(message string, opt ...FuseOpt) error
	R403Forbidden(message string, opt ...FuseOpt) error
	R404NotFound(message string, opt ...FuseOpt) error
	R406NotAcceptable(message string, opt ...FuseOpt) error
	R428PreconditionRequired(message string, opt ...FuseOpt) error
	R500InternalServerError(err error, opt ...FuseOpt) error
}

type FuseOpt struct {
	code     int
	Status   string
	Message  string
	Address  string
	Error    error
	MetaData interface{}
	Data     interface{}
	NewMeta  map[string]interface{}
}

type FuseRouter interface {
	Single(path string, handlers ...func(sc FuseContext) error)
	Group(endpoints map[string][]func(sc FuseContext) error)
}

type FuseDefaultMessage struct {
	R400BadRequest           string
	R401Unauthorized         string
	R403Forbidden            string
	R404NotFound             string
	R406NotAcceptable        string
	R428PreconditionRequired string
	R500InternalServerError  string
}

type FuseDefaultStatus struct {
	R200OK                   string
	R400BadRequest           string
	R401Unauthorized         string
	R403Forbidden            string
	R404NotFound             string
	R406NotAcceptable        string
	R428PreconditionRequired string
	R500InternalServerError  string
}

type srFuseGrpcContext struct {
	ctx      context.Context
	request  *grf.Request
	response *grf.Response
	header   map[string]string
	payload  map[string]*structpb.Value
	params   map[string]string
}

type srFuseContext struct {
	fiberCtx *fiber.Ctx
	grpcCtx  *srFuseGrpcContext
	path     string

	header        map[string]string
	multipartFile map[string][]*multipart.FileHeader

	authX     interface{}
	authY     interface{}
	authZ     interface{}
	isAuthSet bool
}

type srFuseGrpc struct {
	grf.UnimplementedRestfulServiceServer
	routes map[string][]func(FuseContext) error
}

type srFuseRouter struct {
	fiberApp *fiber.App
	fuseGrpc *srFuseGrpc
}

type srConnectFuseGRPC struct {
	address     string
	connMutex   sync.Mutex
	clientMutex sync.Mutex
	conn        *grpc.ClientConn
	client      grf.RestfulServiceClient
}

type ClientFuseGRPC interface {
	Restful(path string, header map[string]string, payload map[string]interface{}, params map[string]string) (*grf.Response, error)
}

type srClientFuseGRPC struct {
	fuseGrpc *srConnectFuseGRPC
}
