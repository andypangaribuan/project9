// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: proto/clog.proto

package clog

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CLogServiceClient is the client API for CLogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CLogServiceClient interface {
	ServiceLog(ctx context.Context, in *RequestServiceLog, opts ...grpc.CallOption) (*Response, error)
	InfoLog(ctx context.Context, in *RequestInfoLog, opts ...grpc.CallOption) (*Response, error)
}

type cLogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCLogServiceClient(cc grpc.ClientConnInterface) CLogServiceClient {
	return &cLogServiceClient{cc}
}

func (c *cLogServiceClient) ServiceLog(ctx context.Context, in *RequestServiceLog, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/clog.CLogService/ServiceLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cLogServiceClient) InfoLog(ctx context.Context, in *RequestInfoLog, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/clog.CLogService/InfoLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CLogServiceServer is the server API for CLogService service.
// All implementations must embed UnimplementedCLogServiceServer
// for forward compatibility
type CLogServiceServer interface {
	ServiceLog(context.Context, *RequestServiceLog) (*Response, error)
	InfoLog(context.Context, *RequestInfoLog) (*Response, error)
	mustEmbedUnimplementedCLogServiceServer()
}

// UnimplementedCLogServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCLogServiceServer struct {
}

func (UnimplementedCLogServiceServer) ServiceLog(context.Context, *RequestServiceLog) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceLog not implemented")
}
func (UnimplementedCLogServiceServer) InfoLog(context.Context, *RequestInfoLog) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InfoLog not implemented")
}
func (UnimplementedCLogServiceServer) mustEmbedUnimplementedCLogServiceServer() {}

// UnsafeCLogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CLogServiceServer will
// result in compilation errors.
type UnsafeCLogServiceServer interface {
	mustEmbedUnimplementedCLogServiceServer()
}

func RegisterCLogServiceServer(s grpc.ServiceRegistrar, srv CLogServiceServer) {
	s.RegisterService(&CLogService_ServiceDesc, srv)
}

func _CLogService_ServiceLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestServiceLog)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CLogServiceServer).ServiceLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clog.CLogService/ServiceLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CLogServiceServer).ServiceLog(ctx, req.(*RequestServiceLog))
	}
	return interceptor(ctx, in, info, handler)
}

func _CLogService_InfoLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInfoLog)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CLogServiceServer).InfoLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clog.CLogService/InfoLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CLogServiceServer).InfoLog(ctx, req.(*RequestInfoLog))
	}
	return interceptor(ctx, in, info, handler)
}

// CLogService_ServiceDesc is the grpc.ServiceDesc for CLogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CLogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clog.CLogService",
	HandlerType: (*CLogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServiceLog",
			Handler:    _CLogService_ServiceLog_Handler,
		},
		{
			MethodName: "InfoLog",
			Handler:    _CLogService_InfoLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/clog.proto",
}
