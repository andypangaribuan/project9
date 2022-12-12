/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package svc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/andypangaribuan/project9/p9"
	"github.com/andypangaribuan/project9/proto/clog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type srCLog struct {
	address string
	mutex   sync.RWMutex
	conn    *grpc.ClientConn
	client  clog.CLogServiceClient
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

func (slf *srCLog) getClient() (clog.CLogServiceClient, error) {
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

	slf.client = clog.NewCLogServiceClient(conn)
	return nil
}

func (slf *srCLog) Info(val CLogRequestInfo) (status string, message string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := &clog.RequestInfoLog{
		Uid:       val.Uid,
		SvcName:   val.SvcName,
		Message:   val.Message,
		Severity:  val.Severity,
		Path:      val.Path,
		Function:  val.Function,
		CreatedAt: p9.Conv.Time.ToStrRFC3339MilliSecond(val.CreatedAt),
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

	req := &clog.RequestServiceLog{
		Uid:       val.Uid,
		SvcName:   val.SvcName,
		Severity:  val.Severity,
		Path:      val.Path,
		Function:  val.Function,
		StartAt:   p9.Conv.Time.ToStrRFC3339MilliSecond(val.StartAt),
		FinishAt:  p9.Conv.Time.ToStrRFC3339MilliSecond(val.FinishAt),
		CreatedAt: p9.Conv.Time.ToStrRFC3339MilliSecond(val.CreatedAt),
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
	if val.ReqPar != nil {
		req.ReqPar = &wrapperspb.StringValue{Value: *val.ReqPar}
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
