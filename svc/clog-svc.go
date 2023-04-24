/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

/* spell-checker: disable */
package svc

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/andypangaribuan/project9/p9"
	"github.com/andypangaribuan/project9/proto/clog-svc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type srCLog struct {
	address  string
	usingTLS bool
	mutex    sync.RWMutex
	conn     *grpc.ClientConn
	client   clog_svc.CLogServiceClient
}

func (slf *srCLog) getConnection() (*grpc.ClientConn, error) {
	if slf.conn == nil {
		err := slf.buildConnection()
		if err != nil {
			return nil, err
		}
	}
	return slf.conn, nil
}

func (slf *srCLog) buildConnection() error {
	if slf.address == "" {
		return errors.New("clog address is empty")
	}

	slf.mutex.Lock()
	defer slf.mutex.Unlock()
	if slf.conn != nil {
		return nil
	}

	var (
		conn *grpc.ClientConn
		err  error
	)

	resolver.SetDefaultScheme("dns")
	const grpcServiceConfig = `{"loadBalancingPolicy":"round_robin"}`

	if slf.usingTLS {
		creds := credentials.NewTLS(&tls.Config{})
		conn, err = grpc.Dial(slf.address, grpc.WithTransportCredentials(creds), grpc.WithDefaultServiceConfig(grpcServiceConfig))
	} else {
		conn, err = grpc.Dial(slf.address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(grpcServiceConfig))
	}

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

func (slf *srCLog) getClient() (clog_svc.CLogServiceClient, error) {
	if slf.client == nil {
		err := slf.buildClient()
		if err != nil {
			return nil, err
		}
	}
	return slf.client, nil
}

func (slf *srCLog) buildClient() error {
	slf.mutex.Lock()
	defer slf.mutex.Unlock()
	if slf.client != nil {
		return nil
	}

	conn, err := slf.getConnection()
	if err != nil {
		return err
	}

	slf.client = clog_svc.NewCLogServiceClient(conn)
	return nil
}

func (slf *srCLog) Info(val CLogRequestInfo) (status string, message string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := &clog_svc.RequestInfoLog{
		Uid:        val.Uid,
		SvcName:    val.SvcName,
		SvcVersion: val.SvcVersion,
		Message:    val.Message,
		Severity:   val.Severity,
		Path:       val.Path,
		Function:   val.Function,
		CreatedAt:  p9.Conv.Time.ToStrRFC3339MilliSecond(val.CreatedAt),
	}

	if val.UserId != nil {
		req.UserId = &wrapperspb.StringValue{Value: *val.UserId}
	}
	if val.PartnerId != nil {
		req.PartnerId = &wrapperspb.StringValue{Value: *val.PartnerId}
	}
	if val.XID != nil {
		req.Xid = &wrapperspb.StringValue{Value: *val.XID}
	}
	if val.SvcParent != nil {
		req.SvcParent = &wrapperspb.StringValue{Value: *val.SvcParent}
	}

	if val.Data != nil {
		req.Data = &wrapperspb.StringValue{Value: *val.Data}
	}

	client, err := slf.getClient()
	if err != nil {
		return "", "", err
	}

	res, err := client.InfoLog(ctx, req)
	if err != nil {
		return "", "", err
	}

	return res.Status, res.Message, nil
}

func (slf *srCLog) Service(val CLogRequestService) (status string, message string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := &clog_svc.RequestServiceLog{
		Uid:        val.Uid,
		SvcName:    val.SvcName,
		SvcVersion: val.SvcVersion,
		Endpoint:   val.Endpoint,
		Version:    val.Version,
		Severity:   val.Severity,
		Path:       val.Path,
		Function:   val.Function,
		ClientIP:   val.ClientIP,
		StartAt:    p9.Conv.Time.ToStrRFC3339MilliSecond(val.StartAt),
		FinishAt:   p9.Conv.Time.ToStrRFC3339MilliSecond(val.FinishAt),
		CreatedAt:  p9.Conv.Time.ToStrRFC3339MilliSecond(val.CreatedAt),
	}

	if val.UserId != nil {
		req.UserId = &wrapperspb.StringValue{Value: *val.UserId}
	}
	if val.PartnerId != nil {
		req.PartnerId = &wrapperspb.StringValue{Value: *val.PartnerId}
	}
	if val.XID != nil {
		req.Xid = &wrapperspb.StringValue{Value: *val.XID}
	}
	if val.SvcParent != nil {
		req.SvcParent = &wrapperspb.StringValue{Value: *val.SvcParent}
	}
	if val.Message != nil {
		req.Message = &wrapperspb.StringValue{Value: *val.Message}
	}
	if val.ReqHeader != nil {
		req.ReqHeader = &wrapperspb.StringValue{Value: *val.ReqHeader}
	}
	if val.ReqBody != nil {
		req.ReqBody = &wrapperspb.StringValue{Value: *val.ReqBody}
	}
	if val.ReqParam != nil {
		req.ReqParam = &wrapperspb.StringValue{Value: *val.ReqParam}
	}
	if val.ResData != nil {
		req.ResData = &wrapperspb.StringValue{Value: *val.ResData}
	}
	if val.ResCode != nil {
		req.ResCode = &wrapperspb.Int32Value{Value: int32(*val.ResCode)}
	}
	if val.Data != nil {
		req.Data = &wrapperspb.StringValue{Value: *val.Data}
	}
	if val.Error != nil {
		req.Error = &wrapperspb.StringValue{Value: *val.Error}
	}
	if val.StackTrace != nil {
		req.StackTrace = &wrapperspb.StringValue{Value: *val.StackTrace}
	}

	client, err := slf.getClient()
	if err != nil {
		return "", "", err
	}

	res, err := client.ServiceLog(ctx, req)
	if err != nil {
		return "", "", err
	}

	return res.Status, res.Message, nil
}

func (slf *srCLog) Dbq(val CLogRequestDbq) (status string, message string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := &clog_svc.RequestDbqLog{
		Uid:        val.Uid,
		SvcName:    val.SvcName,
		SvcVersion: val.SvcVersion,
		SqlQuery:   val.SqlQuery,
		Severity:   val.Severity,
		Path:       val.Path,
		Function:   val.Function,
		StartAt:    p9.Conv.Time.ToStrRFC3339MilliSecond(val.StartAt),
		FinishAt:   p9.Conv.Time.ToStrRFC3339MilliSecond(val.FinishAt),
		CreatedAt:  p9.Conv.Time.ToStrRFC3339MilliSecond(val.CreatedAt),
	}

	if val.UserId != nil {
		req.UserId = &wrapperspb.StringValue{Value: *val.UserId}
	}
	if val.PartnerId != nil {
		req.PartnerId = &wrapperspb.StringValue{Value: *val.PartnerId}
	}
	if val.XID != nil {
		req.Xid = &wrapperspb.StringValue{Value: *val.XID}
	}
	if val.SvcParent != nil {
		req.SvcParent = &wrapperspb.StringValue{Value: *val.SvcParent}
	}
	if val.SqlPars != nil {
		req.SqlPars = &wrapperspb.StringValue{Value: *val.SqlPars}
	}
	if val.Error != nil {
		req.Error = &wrapperspb.StringValue{Value: *val.Error}
	}
	if val.StackTrace != nil {
		req.StackTrace = &wrapperspb.StringValue{Value: *val.StackTrace}
	}

	client, err := slf.getClient()
	if err != nil {
		return "", "", err
	}

	res, err := client.DbqLog(ctx, req)
	if err != nil {
		return "", "", err
	}

	return res.Status, res.Message, nil
}
