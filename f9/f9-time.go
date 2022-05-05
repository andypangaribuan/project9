/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"time"

	"github.com/andypangaribuan/project9/p9"
)

func TimeNow() time.Time {
	dtNowStr := p9.Conv.Time.ToStr(time.Now(), "yyyy-MM-dd HH:mm:ss.SSSSSS") //p9.Conv.Time.ToStrFull(time.Now())
	dtNow, _ := p9.Conv.Time.ToTime("yyyy-MM-dd HH:mm:ss.SSSSSS", dtNowStr)
	return dtNow
}
