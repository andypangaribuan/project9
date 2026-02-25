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
		slf.rawRedisClient = client
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
	ctx := context.Background()
	var lock *redislock.Lock
	var err error
	// Use native retry strategy to respect the tryFor timeout smoothly
	if slf.tryFor != nil && *slf.tryFor > 0 {
		retry := redislock.LinearBackoff(10 * time.Millisecond)
		opts := &redislock.Options{RetryStrategy: retry}
		// Create a bounded context so Obtain automatically gives up after tryFor duration
		tryCtx, cancelTry := context.WithTimeout(ctx, *slf.tryFor)
		defer cancelTry()
		lock, err = slf.redisClient.Obtain(tryCtx, key, *slf.timeout, opts)
	} else {
		// No retry wait time, just try instantly once
		lock, err = slf.redisClient.Obtain(ctx, key, *slf.timeout, nil)
	}
	if err != nil {
		return &srXLock{}, errors.WithStack(err)
	}
	renewCtx, cancel := context.WithCancel(ctx)
	doneChan := make(chan struct{})
	go func() {
		failedToRefresh := 0
		// Refresh when TTL is halfway done
		ticker := time.NewTicker(*slf.timeout / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := lock.Refresh(ctx, *slf.timeout, nil)
				if err == nil {
					failedToRefresh = 0
				} else {
					failedToRefresh++
					// If we fail 3 consecutive times, the lock is definitely dead
					if failedToRefresh == 3 {
						close(doneChan) // Immediately alert the main application!
						return
					}
				}
			case <-renewCtx.Done():
				return // Normal shutdown (Release was called)
			}
		}
	}()
	return &srXLock{
		lockKey:       key,
		lockType:      "redis",
		ctx:           ctx,
		lock:          lock,
		cancel:        &cancel,
		redisDoneChan: doneChan,
	}, nil
}

func (slf *srLock) getEtcdSession() (*concurrency.Session, error) {
	// 1. Fast path: check current session with a Read-Lock
	mxEtcdLock.RLock()
	session := slf.etcdSession
	mxEtcdLock.RUnlock()

	// If session exists and hasn't expired, return it immediately
	if session != nil {
		select {
		case <-session.Done():
			// Session expired or was closed, we need to create a new one
		default:
			return session, nil
		}
	}

	// 2. Slow path: acquire Write-Lock to create a new session
	mxEtcdLock.Lock()
	defer mxEtcdLock.Unlock()
	// Double check in case another goroutine already created it while we waited for the lock
	if slf.etcdSession != nil {
		select {
		case <-slf.etcdSession.Done():
		default:
			return slf.etcdSession, nil
		}
	}

	// Create a single shared session (A 15-second TTL is great for a shared keep-alive)
	newSession, err := concurrency.NewSession(slf.etcdClient, concurrency.WithTTL(15))
	if err != nil {
		return nil, err
	}

	slf.etcdSession = newSession
	return newSession, nil
}

func (slf *srLock) etcdLock(key string) (abs.XLock, error) {
	// 1. Queue locally first to prevent "shared session key-collision" in etcd
	newChan := make(chan struct{}, 1)
	newChan <- struct{}{}
	v, _ := etcdLocalLocks.LoadOrStore(key, newChan)
	localChan := v.(chan struct{})

	var ctx context.Context
	var cancel context.CancelFunc
	if slf.tryFor != nil && *slf.tryFor > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), *slf.tryFor)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	// Wait for local queue first
	select {
	case <-localChan:
		// We safely hold the local lock! Proceed to etcd.
	case <-ctx.Done():
		return &srXLock{}, errors.WithMessage(ctx.Err(), "failed to lock locally (timeout)")
	}

	// 2. Fetch our safe, shared session
	session, err := slf.getEtcdSession()
	if err != nil {
		localChan <- struct{}{} // Abort, let the next local locker try
		return nil, err
	}

	mtx := concurrency.NewMutex(session, "/locks/"+key)

	// 3. Block and wait natively in etcd
	err = mtx.Lock(ctx)
	if err != nil {
		localChan <- struct{}{} // Abort, let the next local locker try
		return &srXLock{}, errors.WithMessage(err, "failed to lock")
	}

	return &srXLock{
		lockKey:       key,
		lockType:      "etcd",
		ctx:           ctx,
		etcdMtx:       mtx,
		etcdSession:   session,
		etcdLocalChan: localChan,
	}, nil
}

func (slf *srLock) Close() {
	// First gracefully handle Redis...
	if slf.rawRedisClient != nil {
		_ = slf.rawRedisClient.Close()
	}

	// Then elegantly close the shared Etcd Session
	if slf.etcdClient != nil {
		mxEtcdLock.Lock()
		defer mxEtcdLock.Unlock()

		if slf.etcdSession != nil {
			// This explicitly revokes the lease on the etcd server side
			// and causes all locks to instantly drop, allowing new pods to take over.
			_ = slf.etcdSession.Close()
			slf.etcdSession = nil
		}

		_ = slf.etcdClient.Close()
	}
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
		if slf.released || slf.etcdMtx == nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := slf.etcdMtx.Unlock(ctx)
		if err != nil {
			log.Printf("error when release etcd lock, error: %v\n", err)
		}

		// Release the local queue, allowing the next goroutine to grab the lock
		if slf.etcdLocalChan != nil {
			select {
			case slf.etcdLocalChan <- struct{}{}:
			default:
			}
		}

		return
	}
}

func (slf *srXLock) Done() <-chan struct{} {
	if slf.lockType == "redis" && slf.redisDoneChan != nil {
		return slf.redisDoneChan
	}

	if slf.lockType == "etcd" && slf.etcdSession != nil {
		// Return the etcd session's underlying kill-switch channel
		return slf.etcdSession.Done()
	}

	// For redis or invalid locks, return a channel that never closes
	// (or implement a context mechanism for Redis lock drops later)
	return make(chan struct{})
}

func (slf *srXLock) clean() {
	slf.released = true
	slf.ctx = nil
	slf.lock = nil
	slf.cancel = nil
	slf.redisDoneChan = nil
	slf.etcdClient = nil
	slf.etcdSession = nil
	slf.etcdMtx = nil
	slf.etcdLocalChan = nil
}
