/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package lock

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/bsm/redislock"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

func (slf *srLock) Init(address string, tryFor time.Duration, timeout time.Duration, engine ...string) {
	address = strings.TrimSpace(address)
	if address == "" || address == "-" {
		return
	}

	useEngine := "redis"
	if len(engine) > 0 {
		useEngine = engine[0]
	}

	if useEngine == "redis" {
		client := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    address,
		})

		slf.redisClient = redislock.New(client)
		slf.timeout = &timeout
		slf.tryFor = &tryFor
		return
	}

	if useEngine == "etcd" {
		ls := strings.Split(address, ",")
		urls := make([]string, 0)
		for _, v := range ls {
			v = strings.TrimSpace(v)
			if v != "" {
				urls = append(urls, v)
			}
		}

		client, err := etcdclientv3.New(etcdclientv3.Config{
			Endpoints:   urls,
			DialTimeout: 3 * time.Second,
			Logger:      zap.NewNop(),
		})

		if err != nil {
			log.Printf("etcd lock create client is error, error: %+v\n", err)
			return
		}

		slf.etcdClient = client
		slf.timeout = &timeout
		slf.tryFor = &tryFor
		return
	}
}

func (slf *srLock) Lock(key string) (abs.XLock, error) {
	if slf.redisClient != nil {
		return slf.redisLock(key)
	}

	if slf.etcdClient != nil {
		return slf.etcdLock(key)
	}

	return &srXLock{}, nil
}

func (slf *srLock) redisLock(key string) (abs.XLock, error) {
	var (
		startedAt = time.Now()
		ctx       = context.Background()
		lock      *redislock.Lock
	)

	for {
		lk, err := slf.redisClient.Obtain(ctx, key, *slf.timeout, nil)
		if err == nil {
			lock = lk
			break
		}

		if slf.tryFor != nil && time.Since(startedAt) > *slf.tryFor {
			return &srXLock{}, errors.WithStack(err)
		}

		time.Sleep(time.Millisecond * 10)
	}

	renewCtx, cancel := context.WithCancel(ctx)
	go func() {
		failedToRefresh := 0
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := lock.Refresh(ctx, *slf.timeout, nil)
				if err == nil {
					failedToRefresh = 0
				} else {
					failedToRefresh++
					if failedToRefresh == 3 {
						return
					}
				}

			case <-renewCtx.Done():
				return
			}
		}
	}()

	return &srXLock{
		lockKey:  key,
		lockType: "redis",
		ctx:      ctx,
		lock:     lock,
		cancel:   &cancel,
	}, nil
}

func (slf *srLock) etcdLock(key string) (abs.XLock, error) {
	ttl := max(int(*slf.timeout/time.Second), 3)
	session, err := concurrency.NewSession(slf.etcdClient, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}

	var (
		startedAt = time.Now()
		ctx       = context.Background()
		mtx       *concurrency.Mutex
	)

	for {
		mtx = concurrency.NewMutex(session, "/locks/"+key)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = mtx.Lock(ctx)
		if err == nil {
			break
		}

		if slf.tryFor != nil && time.Since(startedAt) > *slf.tryFor {
			return &srXLock{}, errors.WithMessage(err, "failed to lock")
		}

		time.Sleep(time.Millisecond * 10)
	}

	return &srXLock{
		lockKey:     key,
		lockType:    "etcd",
		ctx:         ctx,
		etcdClient:  slf.etcdClient,
		etcdSession: session,
		etcdMtx:     mtx,
	}, nil
}

func (slf *srXLock) Release() {
	if slf.lockType == "redis" {
		defer slf.clean()
		if slf.released || slf.lock == nil {
			return
		}

		if slf.cancel != nil {
			cancel := *slf.cancel
			cancel()
		}

		err := slf.lock.Release(slf.ctx)
		if err != nil {
			// log.Printf("error when release redis lock, error: %v\n", err)
			return
		}

		return
	}

	if slf.lockType == "etcd" {
		defer slf.clean()
		if slf.released || slf.etcdMtx == nil && slf.etcdSession == nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := slf.etcdMtx.Unlock(ctx)
		if err != nil {
			// log.Printf("error when release etcd lock, error: %v\n", err)
			return
		}

		err = slf.etcdSession.Close()
		if err != nil {
			// log.Printf("error when do etcd session close, error: %+v\n", err)
			return
		}

		return
	}
}

func (slf *srXLock) clean() {
	slf.released = true
	slf.ctx = nil
	slf.lock = nil
	slf.cancel = nil
	slf.etcdClient = nil
	slf.etcdSession = nil
	slf.etcdMtx = nil
}
