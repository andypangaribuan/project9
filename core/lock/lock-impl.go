/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package lock

import (
	"context"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/bsm/redislock"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (slf *srLock) Init(address string, tryFor time.Duration, timeout time.Duration) {
	address = strings.TrimSpace(address)
	if address == "" || address == "-" {
		return
	}

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    address,
	})

	slf.client = redislock.New(client)
	slf.timeout = &timeout
	slf.tryFor = &tryFor
}

func (slf *srLock) Lock(key string) (abs.XLock, error) {
	if slf.client == nil {
		return &srXLock{}, nil
	}

	var (
		startedAt = time.Now()
		ctx       = context.Background()
		lock      *redislock.Lock
	)

	for {
		lk, err := slf.client.Obtain(ctx, key, *slf.timeout, nil)
		if err == nil {
			lock = lk
			break
		}

		if slf.tryFor != nil && time.Since(startedAt) > *slf.tryFor {
			return nil, errors.WithStack(err)
		}

		time.Sleep(time.Millisecond * 10)
	}

	return &srXLock{ctx: ctx, lock: lock}, nil
}

func (slf *srXLock) Release() {
	if slf.lock == nil {
		return
	}

	_ = slf.lock.Release(slf.ctx)
}
