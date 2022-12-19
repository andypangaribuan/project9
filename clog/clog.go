/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package clog

import (
	"log"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
	psvc "github.com/andypangaribuan/project9/svc"
)

var svc psvc.CLogSVC
var svcName string

func Init(address, serviceName string) {
	svcName = serviceName
	clogSvc, err := psvc.InitCLogSVC(address)
	if err != nil {
		log.Fatal(err)
	}

	svc = clogSvc
}

func SendInfo(depth int, ins Instance, severity Severity, message string, data *string, onGoroutine bool) {
	execFunc, execPath := p9.Util.GetExecutionInfo(1 + depth)
	req := psvc.CLogRequestInfo{
		Uid:       ins.UID,
		SvcName:   svcName,
		SvcParent: f9.Ternary(ins.SvcParent == "", nil, &ins.SvcParent),
		Message:   message,
		Severity:  severity.String(),
		Path:      execPath,
		Function:  execFunc,
		Data:      data,
		CreatedAt: f9.TimeNow(),
	}

	if onGoroutine {
		go func() {
			svcResult(svc.Info(req))
		}()
	} else {
		svcResult(svc.Info(req))
	}
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
		SvcName:    svcName,
		SvcParent:  f9.Ternary(ins.SvcParent == "", nil, &ins.SvcParent),
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
		StartAt:    ins.StartAt,
		FinishAt:   timeNow,
		CreatedAt:  timeNow,
	}

	if onGoroutine {
		go func() {
			svcResult(svc.Service(req))
		}()
	} else {
		svcResult(svc.Service(req))
	}
}

func svcResult(status string, message string, err error) {
	if err != nil {
		log.Printf("clog error: %+v\n", err)
		return
	}

	if status != "success" {
		log.Printf("clog failed, status: %v, message: %v\n", status, message)
	}
}
