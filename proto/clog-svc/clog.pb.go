//
// Copyright (c) 2022.
// Created by Andy Pangaribuan. All Rights Reserved.
//
// This product is protected by copyright and distributed under
// licenses restricting copying, distribution and decompilation.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: proto/clog.proto

package clog_svc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_clog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_clog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_proto_clog_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type RequestServiceLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid        string                  `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	UserId     *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	PartnerId  *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=partnerId,proto3" json:"partnerId,omitempty"`
	Xid        *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=xid,proto3" json:"xid,omitempty"`
	SvcName    string                  `protobuf:"bytes,5,opt,name=svcName,proto3" json:"svcName,omitempty"`
	SvcVersion string                  `protobuf:"bytes,6,opt,name=svcVersion,proto3" json:"svcVersion,omitempty"`
	SvcParent  *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=svcParent,proto3" json:"svcParent,omitempty"`
	Endpoint   string                  `protobuf:"bytes,8,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Version    string                  `protobuf:"bytes,9,opt,name=version,proto3" json:"version,omitempty"`
	Message    *wrapperspb.StringValue `protobuf:"bytes,10,opt,name=message,proto3" json:"message,omitempty"`
	Severity   string                  `protobuf:"bytes,11,opt,name=severity,proto3" json:"severity,omitempty"`
	Path       string                  `protobuf:"bytes,12,opt,name=path,proto3" json:"path,omitempty"`
	Function   string                  `protobuf:"bytes,13,opt,name=function,proto3" json:"function,omitempty"`
	ReqHeader  *wrapperspb.StringValue `protobuf:"bytes,14,opt,name=reqHeader,proto3" json:"reqHeader,omitempty"`
	ReqBody    *wrapperspb.StringValue `protobuf:"bytes,15,opt,name=reqBody,proto3" json:"reqBody,omitempty"`
	ReqParam   *wrapperspb.StringValue `protobuf:"bytes,16,opt,name=reqParam,proto3" json:"reqParam,omitempty"`
	ResData    *wrapperspb.StringValue `protobuf:"bytes,17,opt,name=resData,proto3" json:"resData,omitempty"`
	ResCode    *wrapperspb.Int32Value  `protobuf:"bytes,18,opt,name=resCode,proto3" json:"resCode,omitempty"`
	Data       *wrapperspb.StringValue `protobuf:"bytes,19,opt,name=data,proto3" json:"data,omitempty"`
	Error      *wrapperspb.StringValue `protobuf:"bytes,20,opt,name=error,proto3" json:"error,omitempty"`
	StackTrace *wrapperspb.StringValue `protobuf:"bytes,21,opt,name=stackTrace,proto3" json:"stackTrace,omitempty"`
	ClientIP   string                  `protobuf:"bytes,22,opt,name=clientIP,proto3" json:"clientIP,omitempty"`
	StartAt    string                  `protobuf:"bytes,23,opt,name=startAt,proto3" json:"startAt,omitempty"`     // RFC3339 MilliSecond
	FinishAt   string                  `protobuf:"bytes,24,opt,name=finishAt,proto3" json:"finishAt,omitempty"`   // RFC3339 MilliSecond
	CreatedAt  string                  `protobuf:"bytes,25,opt,name=createdAt,proto3" json:"createdAt,omitempty"` // RFC3339 MilliSecond
}

func (x *RequestServiceLog) Reset() {
	*x = RequestServiceLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_clog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestServiceLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestServiceLog) ProtoMessage() {}

func (x *RequestServiceLog) ProtoReflect() protoreflect.Message {
	mi := &file_proto_clog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestServiceLog.ProtoReflect.Descriptor instead.
func (*RequestServiceLog) Descriptor() ([]byte, []int) {
	return file_proto_clog_proto_rawDescGZIP(), []int{1}
}

func (x *RequestServiceLog) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RequestServiceLog) GetUserId() *wrapperspb.StringValue {
	if x != nil {
		return x.UserId
	}
	return nil
}

func (x *RequestServiceLog) GetPartnerId() *wrapperspb.StringValue {
	if x != nil {
		return x.PartnerId
	}
	return nil
}

func (x *RequestServiceLog) GetXid() *wrapperspb.StringValue {
	if x != nil {
		return x.Xid
	}
	return nil
}

func (x *RequestServiceLog) GetSvcName() string {
	if x != nil {
		return x.SvcName
	}
	return ""
}

func (x *RequestServiceLog) GetSvcVersion() string {
	if x != nil {
		return x.SvcVersion
	}
	return ""
}

func (x *RequestServiceLog) GetSvcParent() *wrapperspb.StringValue {
	if x != nil {
		return x.SvcParent
	}
	return nil
}

func (x *RequestServiceLog) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *RequestServiceLog) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *RequestServiceLog) GetMessage() *wrapperspb.StringValue {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *RequestServiceLog) GetSeverity() string {
	if x != nil {
		return x.Severity
	}
	return ""
}

func (x *RequestServiceLog) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RequestServiceLog) GetFunction() string {
	if x != nil {
		return x.Function
	}
	return ""
}

func (x *RequestServiceLog) GetReqHeader() *wrapperspb.StringValue {
	if x != nil {
		return x.ReqHeader
	}
	return nil
}

func (x *RequestServiceLog) GetReqBody() *wrapperspb.StringValue {
	if x != nil {
		return x.ReqBody
	}
	return nil
}

func (x *RequestServiceLog) GetReqParam() *wrapperspb.StringValue {
	if x != nil {
		return x.ReqParam
	}
	return nil
}

func (x *RequestServiceLog) GetResData() *wrapperspb.StringValue {
	if x != nil {
		return x.ResData
	}
	return nil
}

func (x *RequestServiceLog) GetResCode() *wrapperspb.Int32Value {
	if x != nil {
		return x.ResCode
	}
	return nil
}

func (x *RequestServiceLog) GetData() *wrapperspb.StringValue {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RequestServiceLog) GetError() *wrapperspb.StringValue {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *RequestServiceLog) GetStackTrace() *wrapperspb.StringValue {
	if x != nil {
		return x.StackTrace
	}
	return nil
}

func (x *RequestServiceLog) GetClientIP() string {
	if x != nil {
		return x.ClientIP
	}
	return ""
}

func (x *RequestServiceLog) GetStartAt() string {
	if x != nil {
		return x.StartAt
	}
	return ""
}

func (x *RequestServiceLog) GetFinishAt() string {
	if x != nil {
		return x.FinishAt
	}
	return ""
}

func (x *RequestServiceLog) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type RequestInfoLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid        string                  `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	UserId     *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	PartnerId  *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=partnerId,proto3" json:"partnerId,omitempty"`
	Xid        *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=xid,proto3" json:"xid,omitempty"`
	SvcName    string                  `protobuf:"bytes,5,opt,name=svcName,proto3" json:"svcName,omitempty"`
	SvcVersion string                  `protobuf:"bytes,6,opt,name=svcVersion,proto3" json:"svcVersion,omitempty"`
	SvcParent  *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=svcParent,proto3" json:"svcParent,omitempty"`
	Message    string                  `protobuf:"bytes,8,opt,name=message,proto3" json:"message,omitempty"`
	Severity   string                  `protobuf:"bytes,9,opt,name=severity,proto3" json:"severity,omitempty"`
	Path       string                  `protobuf:"bytes,10,opt,name=path,proto3" json:"path,omitempty"`
	Function   string                  `protobuf:"bytes,11,opt,name=function,proto3" json:"function,omitempty"`
	Data       *wrapperspb.StringValue `protobuf:"bytes,12,opt,name=data,proto3" json:"data,omitempty"`
	CreatedAt  string                  `protobuf:"bytes,13,opt,name=createdAt,proto3" json:"createdAt,omitempty"` // RFC3339 MilliSecond
}

func (x *RequestInfoLog) Reset() {
	*x = RequestInfoLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_clog_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestInfoLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestInfoLog) ProtoMessage() {}

func (x *RequestInfoLog) ProtoReflect() protoreflect.Message {
	mi := &file_proto_clog_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestInfoLog.ProtoReflect.Descriptor instead.
func (*RequestInfoLog) Descriptor() ([]byte, []int) {
	return file_proto_clog_proto_rawDescGZIP(), []int{2}
}

func (x *RequestInfoLog) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RequestInfoLog) GetUserId() *wrapperspb.StringValue {
	if x != nil {
		return x.UserId
	}
	return nil
}

func (x *RequestInfoLog) GetPartnerId() *wrapperspb.StringValue {
	if x != nil {
		return x.PartnerId
	}
	return nil
}

func (x *RequestInfoLog) GetXid() *wrapperspb.StringValue {
	if x != nil {
		return x.Xid
	}
	return nil
}

func (x *RequestInfoLog) GetSvcName() string {
	if x != nil {
		return x.SvcName
	}
	return ""
}

func (x *RequestInfoLog) GetSvcVersion() string {
	if x != nil {
		return x.SvcVersion
	}
	return ""
}

func (x *RequestInfoLog) GetSvcParent() *wrapperspb.StringValue {
	if x != nil {
		return x.SvcParent
	}
	return nil
}

func (x *RequestInfoLog) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RequestInfoLog) GetSeverity() string {
	if x != nil {
		return x.Severity
	}
	return ""
}

func (x *RequestInfoLog) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RequestInfoLog) GetFunction() string {
	if x != nil {
		return x.Function
	}
	return ""
}

func (x *RequestInfoLog) GetData() *wrapperspb.StringValue {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RequestInfoLog) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type RequestDbqLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid        string                  `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	UserId     *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	PartnerId  *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=partnerId,proto3" json:"partnerId,omitempty"`
	Xid        *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=xid,proto3" json:"xid,omitempty"`
	SvcName    string                  `protobuf:"bytes,5,opt,name=svcName,proto3" json:"svcName,omitempty"`
	SvcVersion string                  `protobuf:"bytes,6,opt,name=svcVersion,proto3" json:"svcVersion,omitempty"`
	SvcParent  *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=svcParent,proto3" json:"svcParent,omitempty"`
	SqlQuery   string                  `protobuf:"bytes,8,opt,name=sqlQuery,proto3" json:"sqlQuery,omitempty"`
	SqlPars    *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=sqlPars,proto3" json:"sqlPars,omitempty"`
	Severity   string                  `protobuf:"bytes,10,opt,name=severity,proto3" json:"severity,omitempty"`
	Path       string                  `protobuf:"bytes,11,opt,name=path,proto3" json:"path,omitempty"`
	Function   string                  `protobuf:"bytes,12,opt,name=function,proto3" json:"function,omitempty"`
	Error      *wrapperspb.StringValue `protobuf:"bytes,13,opt,name=error,proto3" json:"error,omitempty"`
	StackTrace *wrapperspb.StringValue `protobuf:"bytes,14,opt,name=stackTrace,proto3" json:"stackTrace,omitempty"`
	StartAt    string                  `protobuf:"bytes,15,opt,name=startAt,proto3" json:"startAt,omitempty"`     // RFC3339 MilliSecond
	FinishAt   string                  `protobuf:"bytes,16,opt,name=finishAt,proto3" json:"finishAt,omitempty"`   // RFC3339 MilliSecond
	CreatedAt  string                  `protobuf:"bytes,17,opt,name=createdAt,proto3" json:"createdAt,omitempty"` // RFC3339 MilliSecond
}

func (x *RequestDbqLog) Reset() {
	*x = RequestDbqLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_clog_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestDbqLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestDbqLog) ProtoMessage() {}

func (x *RequestDbqLog) ProtoReflect() protoreflect.Message {
	mi := &file_proto_clog_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestDbqLog.ProtoReflect.Descriptor instead.
func (*RequestDbqLog) Descriptor() ([]byte, []int) {
	return file_proto_clog_proto_rawDescGZIP(), []int{3}
}

func (x *RequestDbqLog) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RequestDbqLog) GetUserId() *wrapperspb.StringValue {
	if x != nil {
		return x.UserId
	}
	return nil
}

func (x *RequestDbqLog) GetPartnerId() *wrapperspb.StringValue {
	if x != nil {
		return x.PartnerId
	}
	return nil
}

func (x *RequestDbqLog) GetXid() *wrapperspb.StringValue {
	if x != nil {
		return x.Xid
	}
	return nil
}

func (x *RequestDbqLog) GetSvcName() string {
	if x != nil {
		return x.SvcName
	}
	return ""
}

func (x *RequestDbqLog) GetSvcVersion() string {
	if x != nil {
		return x.SvcVersion
	}
	return ""
}

func (x *RequestDbqLog) GetSvcParent() *wrapperspb.StringValue {
	if x != nil {
		return x.SvcParent
	}
	return nil
}

func (x *RequestDbqLog) GetSqlQuery() string {
	if x != nil {
		return x.SqlQuery
	}
	return ""
}

func (x *RequestDbqLog) GetSqlPars() *wrapperspb.StringValue {
	if x != nil {
		return x.SqlPars
	}
	return nil
}

func (x *RequestDbqLog) GetSeverity() string {
	if x != nil {
		return x.Severity
	}
	return ""
}

func (x *RequestDbqLog) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RequestDbqLog) GetFunction() string {
	if x != nil {
		return x.Function
	}
	return ""
}

func (x *RequestDbqLog) GetError() *wrapperspb.StringValue {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *RequestDbqLog) GetStackTrace() *wrapperspb.StringValue {
	if x != nil {
		return x.StackTrace
	}
	return nil
}

func (x *RequestDbqLog) GetStartAt() string {
	if x != nil {
		return x.StartAt
	}
	return ""
}

func (x *RequestDbqLog) GetFinishAt() string {
	if x != nil {
		return x.FinishAt
	}
	return ""
}

func (x *RequestDbqLog) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

var File_proto_clog_proto protoreflect.FileDescriptor

var file_proto_clog_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x63, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x76, 0x63, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a, 0x08,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xa8, 0x08, 0x0a, 0x11, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x67,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x69, 0x64, 0x12, 0x34, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74,
	0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x6e,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x03, 0x78, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x03, 0x78, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x76, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x76, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x73, 0x76, 0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x76, 0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3a,
	0x0a, 0x09, 0x73, 0x76, 0x63, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x09, 0x73, 0x76, 0x63, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x36, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x72, 0x65, 0x71, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x36, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x0f, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x07, 0x72, 0x65, 0x71, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x38, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x72, 0x65, 0x71, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x12, 0x36, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x18, 0x11, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x07, 0x72, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x35, 0x0a, 0x07, 0x72, 0x65,
	0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e,
	0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x72, 0x65, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x30, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x3c, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x63, 0x6b,
	0x54, 0x72, 0x61, 0x63, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x63, 0x6b,
	0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x50, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x50, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74, 0x18, 0x17, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x74, 0x18, 0x18, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x18, 0x19, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xf0, 0x03, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x6f, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x34, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x3a, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x03,
	0x78, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x03, 0x78, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x76, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x76, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x76, 0x63, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x76, 0x63, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x09, 0x73, 0x76, 0x63, 0x50, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x73, 0x76, 0x63, 0x50, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x9f, 0x05, 0x0a, 0x0d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x44, 0x62, 0x71, 0x4c, 0x6f, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x34, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x3a, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2e,
	0x0a, 0x03, 0x78, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x03, 0x78, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x76, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x76, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x76, 0x63, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x76,
	0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x09, 0x73, 0x76, 0x63, 0x50,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x73, 0x76, 0x63, 0x50, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x71, 0x6c, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x71, 0x6c, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x12, 0x36, 0x0a, 0x07, 0x73, 0x71, 0x6c, 0x50, 0x61, 0x72, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x07, 0x73, 0x71, 0x6c, 0x50, 0x61, 0x72, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x3c, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x63,
	0x6b, 0x54, 0x72, 0x61, 0x63, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x63,
	0x6b, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41,
	0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x74, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x32, 0xbc, 0x01, 0x0a, 0x0b, 0x43,
	0x4c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x12, 0x1b, 0x2e, 0x63, 0x6c, 0x6f, 0x67, 0x5f,
	0x73, 0x76, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4c, 0x6f, 0x67, 0x1a, 0x12, 0x2e, 0x63, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x76, 0x63,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x07, 0x49, 0x6e, 0x66,
	0x6f, 0x4c, 0x6f, 0x67, 0x12, 0x18, 0x2e, 0x63, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x76, 0x63, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x6f, 0x67, 0x1a, 0x12,
	0x2e, 0x63, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x76, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x44, 0x62, 0x71, 0x4c, 0x6f, 0x67, 0x12, 0x17, 0x2e, 0x63,
	0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x76, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44,
	0x62, 0x71, 0x4c, 0x6f, 0x67, 0x1a, 0x12, 0x2e, 0x63, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x76, 0x63,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x50, 0x0a, 0x1b, 0x69, 0x6f, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x68, 0x65,
	0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x42, 0x0f, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x18, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x63, 0x6c, 0x6f,
	0x67, 0x2d, 0x73, 0x76, 0x63, 0xa2, 0x02, 0x03, 0x48, 0x4c, 0x57, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_clog_proto_rawDescOnce sync.Once
	file_proto_clog_proto_rawDescData = file_proto_clog_proto_rawDesc
)

func file_proto_clog_proto_rawDescGZIP() []byte {
	file_proto_clog_proto_rawDescOnce.Do(func() {
		file_proto_clog_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_clog_proto_rawDescData)
	})
	return file_proto_clog_proto_rawDescData
}

var file_proto_clog_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_clog_proto_goTypes = []interface{}{
	(*Response)(nil),               // 0: clog_svc.Response
	(*RequestServiceLog)(nil),      // 1: clog_svc.RequestServiceLog
	(*RequestInfoLog)(nil),         // 2: clog_svc.RequestInfoLog
	(*RequestDbqLog)(nil),          // 3: clog_svc.RequestDbqLog
	(*wrapperspb.StringValue)(nil), // 4: google.protobuf.StringValue
	(*wrapperspb.Int32Value)(nil),  // 5: google.protobuf.Int32Value
}
var file_proto_clog_proto_depIdxs = []int32{
	4,  // 0: clog_svc.RequestServiceLog.userId:type_name -> google.protobuf.StringValue
	4,  // 1: clog_svc.RequestServiceLog.partnerId:type_name -> google.protobuf.StringValue
	4,  // 2: clog_svc.RequestServiceLog.xid:type_name -> google.protobuf.StringValue
	4,  // 3: clog_svc.RequestServiceLog.svcParent:type_name -> google.protobuf.StringValue
	4,  // 4: clog_svc.RequestServiceLog.message:type_name -> google.protobuf.StringValue
	4,  // 5: clog_svc.RequestServiceLog.reqHeader:type_name -> google.protobuf.StringValue
	4,  // 6: clog_svc.RequestServiceLog.reqBody:type_name -> google.protobuf.StringValue
	4,  // 7: clog_svc.RequestServiceLog.reqParam:type_name -> google.protobuf.StringValue
	4,  // 8: clog_svc.RequestServiceLog.resData:type_name -> google.protobuf.StringValue
	5,  // 9: clog_svc.RequestServiceLog.resCode:type_name -> google.protobuf.Int32Value
	4,  // 10: clog_svc.RequestServiceLog.data:type_name -> google.protobuf.StringValue
	4,  // 11: clog_svc.RequestServiceLog.error:type_name -> google.protobuf.StringValue
	4,  // 12: clog_svc.RequestServiceLog.stackTrace:type_name -> google.protobuf.StringValue
	4,  // 13: clog_svc.RequestInfoLog.userId:type_name -> google.protobuf.StringValue
	4,  // 14: clog_svc.RequestInfoLog.partnerId:type_name -> google.protobuf.StringValue
	4,  // 15: clog_svc.RequestInfoLog.xid:type_name -> google.protobuf.StringValue
	4,  // 16: clog_svc.RequestInfoLog.svcParent:type_name -> google.protobuf.StringValue
	4,  // 17: clog_svc.RequestInfoLog.data:type_name -> google.protobuf.StringValue
	4,  // 18: clog_svc.RequestDbqLog.userId:type_name -> google.protobuf.StringValue
	4,  // 19: clog_svc.RequestDbqLog.partnerId:type_name -> google.protobuf.StringValue
	4,  // 20: clog_svc.RequestDbqLog.xid:type_name -> google.protobuf.StringValue
	4,  // 21: clog_svc.RequestDbqLog.svcParent:type_name -> google.protobuf.StringValue
	4,  // 22: clog_svc.RequestDbqLog.sqlPars:type_name -> google.protobuf.StringValue
	4,  // 23: clog_svc.RequestDbqLog.error:type_name -> google.protobuf.StringValue
	4,  // 24: clog_svc.RequestDbqLog.stackTrace:type_name -> google.protobuf.StringValue
	1,  // 25: clog_svc.CLogService.ServiceLog:input_type -> clog_svc.RequestServiceLog
	2,  // 26: clog_svc.CLogService.InfoLog:input_type -> clog_svc.RequestInfoLog
	3,  // 27: clog_svc.CLogService.DbqLog:input_type -> clog_svc.RequestDbqLog
	0,  // 28: clog_svc.CLogService.ServiceLog:output_type -> clog_svc.Response
	0,  // 29: clog_svc.CLogService.InfoLog:output_type -> clog_svc.Response
	0,  // 30: clog_svc.CLogService.DbqLog:output_type -> clog_svc.Response
	28, // [28:31] is the sub-list for method output_type
	25, // [25:28] is the sub-list for method input_type
	25, // [25:25] is the sub-list for extension type_name
	25, // [25:25] is the sub-list for extension extendee
	0,  // [0:25] is the sub-list for field type_name
}

func init() { file_proto_clog_proto_init() }
func file_proto_clog_proto_init() {
	if File_proto_clog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_clog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_clog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestServiceLog); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_clog_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestInfoLog); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_clog_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestDbqLog); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_clog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_clog_proto_goTypes,
		DependencyIndexes: file_proto_clog_proto_depIdxs,
		MessageInfos:      file_proto_clog_proto_msgTypes,
	}.Build()
	File_proto_clog_proto = out.File
	file_proto_clog_proto_rawDesc = nil
	file_proto_clog_proto_goTypes = nil
	file_proto_clog_proto_depIdxs = nil
}
