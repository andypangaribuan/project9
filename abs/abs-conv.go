/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "time"

type ConvTime interface {
	ToStrDate(tm time.Time) string
	ToStrFull(tm time.Time) string
}
