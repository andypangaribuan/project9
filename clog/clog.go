/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package clog

import (
	"log"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
	psvc "github.com/andypangaribuan/project9/svc"
)

var (
	svc        psvc.CLogSVC
	svcName    string
	svcVersion string
	retryCount int           = 1
	retryAfter time.Duration = time.Second
)

type srSvcResultHandler struct {
	status  string
	message string
}

func Init(address, serviceName, serviceVersion string, retries int, retryDelay time.Duration, usingClientLoadBalancing bool) {
	svcName = serviceName
	svcVersion = serviceVersion
	if retries >= 0 {
		retryCount = retries
		retryAfter = retryDelay
	}

	usingTLS := false
	arr := strings.Split(address, ":")
	if len(arr) > 1 {
		if arr[len(arr)-1] == "443" {
			usingTLS = true
		}
	}

	clogSvc, err := psvc.InitCLogSVC(address, usingTLS, usingClientLoadBalancing)
	if err != nil {
		log.Fatal(err)
	}

	svc = clogSvc
}

func SendInfo(depth int, ins Instance, severity Severity, message string, data *string, onGoroutine bool) {
	execFunc, execPath := p9.Util.GetExecutionInfo(1 + depth)
	req := psvc.CLogRequestInfo{
		Uid:        ins.UID,
		UserId:     ins.UserId,
		PartnerId:  ins.PartnerId,
		XID:        ins.XID,
		SvcName:    svcName,
		SvcVersion: svcVersion,
		SvcParent:  f9.Ternary(ins.SvcParent == "", nil, &ins.SvcParent),
		Message:    message,
		Severity:   severity.String(),
		Path:       execPath,
		Function:   execFunc,
		Data:       data,
		CreatedAt:  f9.TimeNow(),
	}

	send(&req, nil, nil, onGoroutine)
}

func SendService(depth int, ins Instance, severity Severity, m SendServiceModel, onGoroutine bool) {
	var (
		timeNow  = f9.TimeNow()
		execFunc string
		execPath string
	)

	if m.ExecFunc != nil && m.ExecPath != nil {
		execFunc = *m.ExecFunc
		execPath = *m.ExecPath
	} else {
		execFunc, execPath = p9.Util.GetExecutionInfo(1 + depth)
	}

	req := psvc.CLogRequestService{
		Uid:        ins.UID,
		UserId:     ins.UserId,
		PartnerId:  ins.PartnerId,
		XID:        ins.XID,
		SvcName:    svcName,
		SvcVersion: svcVersion,
		SvcParent:  f9.Ternary(ins.SvcParent == "", nil, &ins.SvcParent),
		Endpoint:   m.Endpoint,
		Version:    ins.EndpointVersion,
		Message:    m.Message,
		Severity:   severity.String(),
		Path:       execPath,
		Function:   execFunc,
		ReqHeader:  m.ReqHeader,
		ReqBody:    m.ReqBody,
		ReqParam:   m.ReqParam,
		ResData:    m.ResData,
		ResCode:    m.ResCode,
		Data:       m.Data,
		Error:      m.Error,
		StackTrace: m.StackTrace,
		ClientIP:   m.ClientIP,
		StartAt:    ins.StartAt,
		FinishAt:   timeNow,
		CreatedAt:  timeNow,
	}

	send(nil, &req, nil, onGoroutine)
}

func SendDbq(depth int, ins Instance, severity Severity, m SendDbqModel, onGoroutine bool) {
	var (
		timeNow  = f9.TimeNow()
		execFunc string
		execPath string
	)

	if m.ExecFunc != nil && m.ExecPath != nil {
		execFunc = *m.ExecFunc
		execPath = *m.ExecPath
	} else {
		execFunc, execPath = p9.Util.GetExecutionInfo(1 + depth)
	}

	req := psvc.CLogRequestDbq{
		Uid:        ins.UID,
		UserId:     ins.UserId,
		PartnerId:  ins.PartnerId,
		XID:        ins.XID,
		SvcName:    svcName,
		SvcVersion: svcVersion,
		SvcParent:  f9.Ternary(ins.SvcParent == "", nil, &ins.SvcParent),
		SqlQuery:   m.SqlQuery,
		SqlPars:    m.SqlPars,
		Severity:   severity.String(),
		Path:       execPath,
		Function:   execFunc,
		Error:      m.Error,
		StackTrace: m.StackTrace,
		StartAt:    m.StartAt,
		FinishAt:   timeNow,
		CreatedAt:  timeNow,
	}

	send(nil, nil, &req, onGoroutine)
}

func send(infoData *psvc.CLogRequestInfo, serviceData *psvc.CLogRequestService, dbqData *psvc.CLogRequestDbq, onGoroutine bool) {
	if onGoroutine {
		go func() {
			send(infoData, serviceData, dbqData, false)
		}()
		return
	}

	exec := func() (*srSvcResultHandler, error) {
		var (
			status  string
			message string
			err     error
		)

		switch {
		case infoData != nil:
			status, message, err = svc.Info(*infoData)

		case serviceData != nil:
			status, message, err = svc.Service(*serviceData)

		case dbqData != nil:
			status, message, err = svc.Dbq(*dbqData)
		}

		res := &srSvcResultHandler{status: status, message: message}
		return res, err
	}

	res, err := f9.Retry(exec, retryCount, retryAfter)
	svcResult(res, err)
}

func svcResult(res *srSvcResultHandler, err error) {
	if err != nil {
		log.Printf("clog error: %+v\n", err)
		return
	}

	if res == nil {
		log.Printf("unavailable on retry action")
		return
	}

	if res.status != "success" {
		log.Printf("clog failed, status: %v, message: %v\n", res.status, res.message)
	}
}
