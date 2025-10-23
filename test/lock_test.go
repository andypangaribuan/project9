/*
 * Copyright (c) 2025.
 * Created by Andy Pangaribuan (iam.pangaribuan@gmail.com)
 * https://github.com/apangaribuan
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 * All Rights Reserved.
 */

package test

import (
	"log"
	"testing"
	"time"

	"github.com/andypangaribuan/project9/p9"
	"github.com/stretchr/testify/assert"
)

func TestRedisDisLock(t *testing.T) {
	address := "localhost:6379"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "redis")

	lock, err := p9.Lock.Lock(key)
	defer lock.Release()

	assert.Nil(t, err)
	if err != nil {
		return
	}

	time.Sleep(time.Second * 30)
}

func TestRedisDisLockError(t *testing.T) {
	address := "localhost:6379"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "redis")

	lock, err := p9.Lock.Lock(key)
	defer lock.Release()

	assert.Nil(t, err)
	if err != nil {
		return
	}

	time.Sleep(time.Second * 10)
}

func TestEtcdDisLock(t *testing.T) {
	address := "localhost:2790, localhost:2791, localhost:2792"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "etcd")

	lock, err := p9.Lock.Lock(key)
	defer lock.Release()

	assert.Nil(t, err)
	if err != nil {
		return
	}

	time.Sleep(time.Second * 30)
}

func TestEtcdDisLockError(t *testing.T) {
	go etcdDisLockErrorLogic("c-1")

	time.Sleep(time.Second * 3)
	go etcdDisLockErrorLogic("c-2")

	time.Sleep(time.Second * 1)
	go etcdDisLockErrorLogic("c-3")

	time.Sleep(time.Second * 1)
	go etcdDisLockErrorLogic("c-4")

	time.Sleep(time.Second * 4)
	go etcdDisLockErrorLogic("c-5")
	go etcdDisLockErrorLogic("c-6")
	go etcdDisLockErrorLogic("c-7")

	time.Sleep(time.Second * 30)
	go etcdDisLockErrorLogic("c-10")

	time.Sleep(time.Second * 90)
}

func etcdDisLockErrorLogic(clientName string) {
	address := "localhost:2790, localhost:2791, localhost:2792"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "etcd")

	lock, err := p9.Lock.Lock(key)
	defer lock.Release()

	if err != nil {
		log.Printf("cannot get the lock, client: %v\n", clientName)
		return
	}

	log.Println("locked, client: " + clientName)
	time.Sleep(time.Second * 5)
}
