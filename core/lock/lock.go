/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package lock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

type srLock struct {
	client  *redislock.Client
	timeout *time.Duration
	tryFor  *time.Duration
}

type srXLock struct {
	ctx  context.Context
	lock *redislock.Lock
}

func Create() *srLock {
	return &srLock{}
}
