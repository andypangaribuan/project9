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
	UID       string
	SvcParent string
	StartAt   time.Time
}

type SendServiceModel struct {
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
}
