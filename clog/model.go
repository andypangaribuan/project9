/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package clog

import (
	"time"
)

type Instance struct {
	UID                 string
	SvcParent           string
	UserId              *string
	PartnerId           *string
	XID                 *string
	EndpointVersion     string
	StartAt             time.Time
	ReqBodySaveExcluded []string
}

type SendServiceModel struct {
	Endpoint   string
	ExecFunc   *string
	ExecPath   *string
	Message    *string
	ReqHeader  *string
	ReqBody    *string
	ReqParam   *string
	ResData    *string
	ResCode    *int
	Data       *string
	Error      *string
	StackTrace *string
	ClientIP   string
}

type SendDbqModel struct {
	StartAt    time.Time
	ExecFunc   *string
	ExecPath   *string
	SqlQuery   string
	SqlPars    *string
	Error      *string
	StackTrace *string
}
