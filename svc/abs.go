/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package svc

type CLogSVC interface {
	Info(val CLogRequestInfo) (status string, message string, err error)
	Service(val CLogRequestService) (status string, message string, err error)
}
