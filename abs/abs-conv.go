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
	ToTime(layout string, value string) (tm time.Time, err error)
}
