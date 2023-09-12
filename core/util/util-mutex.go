/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/viney-shih/go-lock"
)

const min int64 = 10
const max int64 = 500

func (*srUtil) NewMutex(name string) abs.UtilMutex {
	return &srMutex{
		mux:  lock.NewChanMutex(),
		name: name,
	}
}

func (slf *srMutex) Sleep(duration ...time.Duration) {
	if len(duration) > 0 {
		time.Sleep(duration[0])
	} else {
		x := rand.Int63n(max-min) + min
		time.Sleep(time.Microsecond * time.Duration(x))
	}
}

func (slf *srMutex) Unlock() {
	slf.mux.Unlock()
}

func (slf *srMutex) Lock(timeout ...time.Duration) (isTimeout bool) {
	if len(timeout) > 0 && timeout[0] > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout[0])
		defer cancel()

		locked := slf.mux.TryLockWithContext(ctx)
		if !locked {
			isTimeout = true
		}

		return
	}

	slf.mux.Lock()
	return
}

func (slf *srMutex) Exec(timeout *time.Duration, fn func()) (executed bool, panicErr error) {
	var (
		isUnlock        = false
		timeoutDuration = make([]time.Duration, 0)
	)

	if timeout != nil {
		timeoutDuration = append(timeoutDuration, *timeout)
	}

	isTimeout := slf.Lock(timeoutDuration...)
	if isTimeout {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr = fmt.Errorf("panic error: %+v", r)
		}

		if !isUnlock {
			slf.Unlock()
		}
	}()

	fn()
	executed = true

	slf.Unlock()
	isUnlock = true

	return
}

func (slf *srMutex) FExec(timeoutLock *time.Duration, timeoutFunc time.Duration, fn func()) (executed bool, isTimeout bool, panicErr error) {
	executed, panicErr = slf.Exec(timeoutLock, func() {
		isTimeout, panicErr = slf.Func(timeoutFunc, fn)
	})

	return
}

func (slf *srMutex) Func(timeout time.Duration, fn func()) (isTimeout bool, panicErr error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("panic error: %+v", r)
			}
			cancel()
		}()

		fn()
	}(ctx)

	<-ctx.Done()
	if panicErr != nil {
		return
	}

	switch ctx.Err() {
	case context.DeadlineExceeded:
		isTimeout = true
	}

	return
}
