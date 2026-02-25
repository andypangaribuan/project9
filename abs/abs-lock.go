/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "time"

type Lock interface {
	Init(address string, tryFor time.Duration, timeout time.Duration, engine ...string)
	Lock(key string) (XLock, error)
	Close()
}

type XLock interface {
	Release()
	Done() <-chan struct{}
}
