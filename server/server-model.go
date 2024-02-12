/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

/* spell-checker: disable */
package server

import (
	"context"
	"mime/multipart"
	"net"
	"sync"

	"github.com/andypangaribuan/project9/clog"
	"github.com/andypangaribuan/project9/server/proto/gen/grf"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type FuseContext interface {
	URI() string
	Request() FuseContextRequest
	Params(key string, defaultValue ...string) string
	Query(key string, defaultValue ...string) string
	Parser(logc *clog.Instance, header, body interface{}) (bool, error)
	GetMultipartFiles() map[string][]*multipart.FileHeader
	ClientIP() string
	Path() string
	Method() string

	AuthX() interface{}
	AuthY() interface{}
	AuthZ() interface{}
	SetAuth(authX, authY, authZ interface{})

	SetCLog(logc *clog.Instance)
	GetCLog() *clog.Instance

	GetResCode() (int, string)
	GetResObject() interface{}
	SetResponse(code int, subCode string, obj interface{}) error
	SR200OK(subCode string, obj interface{}) error
	SR400BadRequest(subCode string, obj interface{}) error
	SR401Unauthorized(subCode string, obj interface{}) error
	SR403Forbidden(subCode string, obj interface{}) error
	SR404NotFound(subCode string, obj interface{}) error
	SR406NotAcceptable(subCode string, obj interface{}) error
	SR428PreconditionRequired(subCode string, obj interface{}) error
	SR500InternalServerError(subCode string, obj interface{}) error

	RString(logc *clog.Instance, code int, data string) error
	RJson(logc *clog.Instance, code int, data interface{}, opt ...FuseOpt) error
	RJsonRaw(logc *clog.Instance, code int, data []byte) error

	Redirect(logc *clog.Instance, code int, url string) error
	RXXX(logc *clog.Instance, code int, status string, data interface{}, opt ...FuseOpt) error
	R200OK(logc *clog.Instance, data interface{}, opt ...FuseOpt) error
	R400BadRequest(logc *clog.Instance, message string, opt ...FuseOpt) error
	R401Unauthorized(logc *clog.Instance, message string, opt ...FuseOpt) error
	R403Forbidden(logc *clog.Instance, message string, opt ...FuseOpt) error
	R404NotFound(logc *clog.Instance, message string, opt ...FuseOpt) error
	R406NotAcceptable(logc *clog.Instance, message string, opt ...FuseOpt) error
	R428PreconditionRequired(logc *clog.Instance, message string, opt ...FuseOpt) error
	R500InternalServerError(logc *clog.Instance, err error, opt ...FuseOpt) error

	SetSendResponse(send bool)
	GetUnSendResponse() *FuseResponse
	GetUnSendResponseOpt() []FuseOpt
}

type FuseContextRequest interface {
	Body() []byte
	Header() map[string]string
	SetHeader(key, value string)
}

type FuseOpt struct {
	code               int
	SubCode            string
	Status             string
	Message            string
	Address            string
	Error              error
	MetaData           interface{}
	Data               interface{}
	NewMeta            map[string]interface{}
	NewHeader          map[string]interface{}
	LogMessage         string
	LogData            string
	LogDepthAdd        int
	JsonCustomOutput   bool
	StringCustomOutput bool
}

type FuseRouter interface {
	Unrouted(handler func(ctx FuseContext, method, path, url string) error)
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

type FuseResponse struct {
	Meta FuseResponseMeta `json:"meta"`
	Data interface{}      `json:"data,omitempty"`
}

type FuseResponseMeta struct {
	Code    int         `json:"code"`
	SubCode string      `json:"subcode,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Address string      `json:"address,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type srFuseGrpcContext struct {
	ctx      context.Context
	request  *grf.Request
	response *grf.Response
	header   map[string]string
	payload  map[string]*structpb.Value
	params   map[string]string
	queries  map[string]string
}

type srFuseContext struct {
	fiberCtx *fiber.Ctx
	grpcCtx  *srFuseGrpcContext
	reqCtx   *srFuseContextRequest
	path     string
	clientIP string

	header        map[string]string
	multipartFile map[string][]*multipart.FileHeader
	jsonBody      *string

	authX     interface{}
	authY     interface{}
	authZ     interface{}
	isAuthSet bool

	resCode    int
	resSubCode string
	resObject  interface{}

	logc              *clog.Instance
	sendResponse      bool
	unSendResponse    *FuseResponse
	unSendResponseOpt []FuseOpt
}

type srFuseContextRequest struct {
	fuseCtx *srFuseContext
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

type srClientIP struct {
	cidrs                       []*net.IPNet
	xOriginalForwardedForHeader string
	xForwardedForHeader         string
	xForwardedHeader            string
	forwardedForHeader          string
	forwardedHeader             string
	xClientIPHeader             string
	xRealIPHeader               string
	cfConnectingIPHeader        string
	fastlyClientIPHeader        string
	trueClientIPHeader          string
}
