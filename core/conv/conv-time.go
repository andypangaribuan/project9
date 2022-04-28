/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

import "time"

const (
	layoutTimeDate = "2006-01-02"
	layoutTimeFull = "2006-01-02 15:04:05"
)

func (*srTime) ToStrDate(tm time.Time) string {
	return tm.Format(layoutTimeDate)
}

func (*srTime) ToStrFull(tm time.Time) string {
	return tm.Format(layoutTimeFull)
}
