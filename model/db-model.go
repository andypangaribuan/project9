/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package model

import "time"

type DbConfig struct {
	MaxLifeTimeConnection time.Duration
	MaxIdleTimeConnection time.Duration
	MaxIdleConnection     int
	MaxOpenConnection     int
}

type DbUnsafeSelectError struct {
	LogType    string
	SqlQuery   string
	SqlPars    []interface{}
	LogMessage *string
	LogTrace   *string
}

type DbExec struct {
	Host string
}
