/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

type Conv interface {
	AnyToMap(obj interface{}) (map[string]interface{}, error)
}

type ConvTime interface {
	ToStr(tm time.Time, format string) string
	ToStrDate(tm time.Time) string
	ToStrFull(tm time.Time) string
	ToStrMillis(tm time.Time) string
	ToStrMicro(tm time.Time) string

	ToStrRFC3339(tm time.Time) string
	ToStrRFC3339MilliSecond(tm time.Time) string
	ToStrRFC3339MicroSecond(tm time.Time) string

	ToTime(layout string, value string) (tm time.Time, err error)
	ToTimeDate(value string) (time.Time, error)
	ToTimeFull(value string) (time.Time, error)
	ToTimeMillis(value string) (time.Time, error)
	ToTimeMicro(value string) (time.Time, error)

	GetTimeZone(tm time.Time) string
	RemovePart(tm time.Time, part ...string) time.Time
}

type ConvProto interface {
	AnyToProtoStruct(sr any) (*structpb.Struct, error)
}
