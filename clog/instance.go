/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package clog

import (
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

func New() *Instance {
	return &Instance{
		UID:     p9.Util.GetXID25(),
		StartAt: f9.TimeNow(),
	}
}

func (slf *Instance) SetId(userId, partnerId, xid *string) {
	slf.UserId = userId
	slf.PartnerId = partnerId
	slf.XID = xid
}

func (slf *Instance) Info(message string, data ...string) {
	var logData *string
	if len(data) > 0 {
		logData = &data[0]
	}

	SendInfo(1, *slf, Info, message, logData, true)
}

func (slf *Instance) Warning(message string, data ...string) {
	var logData *string
	if len(data) > 0 {
		logData = &data[0]
	}

	SendInfo(1, *slf, Warning, message, logData, true)
}

func (slf *Instance) Error(message string, data ...string) {
	var logData *string
	if len(data) > 0 {
		logData = &data[0]
	}

	SendInfo(1, *slf, Error, message, logData, true)
}
