/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/viney-shih/go-lock"
)

func (slf *srUtil) NewNMutex(max int) abs.UtilNMutex {
	return &srNMutex{
		sr:   slf,
		mux:  lock.NewChanMutex(),
		max:  max,
		keys: make(map[string]interface{}, 0),
	}
}

func (slf *srNMutex) getTimeout() time.Duration {
	return slf.sr.GetRandomDuration(1, 500, time.Microsecond)
}

func (slf *srNMutex) sleep() {
	time.Sleep(slf.getTimeout())
}

func (slf *srNMutex) Lock(key string, totalTry ...int) (locked bool) {
	n := 1
	if len(totalTry) > 0 {
		n = totalTry[0]
		if n < 1 {
			n = 1
		}
	}

	for i := 0; i < n; i++ {
		locked = slf.doLock(key)
		if locked {
			break
		}
	}

	return
}

func (slf *srNMutex) doLock(key string) (locked bool) {
	locked = slf.mux.TryLockWithTimeout(slf.getTimeout())
	if !locked {
		return
	}
	defer slf.mux.Unlock()

	if len(slf.keys) >= slf.max {
		slf.sleep()
		return false
	}

	slf.keys[key] = nil
	return true
}

func (slf *srNMutex) Unlock(key string) {
	slf.mux.Lock()
	defer slf.mux.Unlock()
	delete(slf.keys, key)
}
