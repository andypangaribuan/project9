/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package svc

import "time"

type CLogRequestInfo struct {
	Uid       string
	SvcName   string
	SvcParent *string
	Message   string
	Severity  string
	Path      string
	Function  string
	Data      *string
	CreatedAt time.Time
}

type CLogRequestService struct {
	Uid        string
	SvcName    string
	SvcParent  *string
	Message    *string
	Severity   string
	Path       string
	Function   string
	ReqHeader  *string
	ReqBody    *string
	ReqPar     *string
	ResData    *string
	ResCode    *int
	Data       *string
	Error      *string
	StackTrace *string
	StartAt    time.Time
	FinishAt   time.Time
	CreatedAt  time.Time
}
