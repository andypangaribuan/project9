/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "time"

type ConvTime interface {
	ToStr(tm time.Time, format string) string
	ToStrDate(tm time.Time) string
	ToStrFull(tm time.Time) string
	ToStrMillis(tm time.Time) string
	ToStrMicro(tm time.Time) string

	ToTime(layout string, value string) (tm time.Time, err error)
	ToTimeDate(value string) (time.Time, error)
	ToTimeFull(value string) (time.Time, error)
	ToTimeMillis(value string) (time.Time, error)
	ToTimeMicro(value string) (time.Time, error)
}
