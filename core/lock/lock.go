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
	"github.com/redis/go-redis/v9"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type srLock struct {
	redisClient    *redislock.Client
	rawRedisClient *redis.Client
	etcdClient     *etcdclientv3.Client
	etcdSession    *concurrency.Session
	timeout        *time.Duration
	tryFor         *time.Duration
}

type srXLock struct {
	released      bool
	lockKey       string
	lockType      string
	ctx           context.Context
	lock          *redislock.Lock
	cancel        *context.CancelFunc
	redisDoneChan chan struct{}
	etcdClient    *etcdclientv3.Client
	etcdSession   *concurrency.Session
	etcdMtx       *concurrency.Mutex
	etcdLocalChan chan struct{}
}

func Create() *srLock {
	return &srLock{}
}
