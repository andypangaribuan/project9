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
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/andypangaribuan/project9/p9"
	"github.com/stretchr/testify/assert"
)

func TestRedisDisLock(t *testing.T) {
	address := "localhost:2379"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "redis")

	lock, err := p9.Lock.Lock(key)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	defer lock.Release()

	time.Sleep(time.Second * 30)
}

func TestRedisDisLockError(t *testing.T) {
	address := "localhost:2379"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "redis")

	lock, err := p9.Lock.Lock(key)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	defer lock.Release()

	time.Sleep(time.Second * 10)
}

func TestEtcdDisLock(t *testing.T) {
	// address := "127.0.0.1:2379"
	address := "127.0.0.1:2379, 127.0.0.1:2380, 127.0.0.1:50686"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "etcd")

	fmt.Println("try to lock")
	lock, err := p9.Lock.Lock(key)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	defer lock.Release()

	fmt.Println("lock acquired")
	time.Sleep(time.Second * 3)
	fmt.Println("unlock")
}

func TestEtcdDisLockLoop(t *testing.T) {
	// cd test && go test -v -run ^TestEtcdDisLockLoop$
	// address := "127.0.0.1:2379"
	address := "127.0.0.1:2379, 127.0.0.1:2380, 127.0.0.1:50686"
	p9.Lock.Init(address, time.Second*5, time.Second*3, "etcd")

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			key := fmt.Sprintf("lock-%v", p9.Util.GetRandomNumber(1, 3))
			log.Printf("try to lock, key: %v\n", key)
			lock, err := p9.Lock.Lock(key)
			if err != nil {
				log.Printf("cannot get the lock, key: %v, error: %v\n", key, err)
				return
			}

			defer lock.Release()
			num := p9.Util.GetRandomNumber(1, 10)
			log.Printf("lock acquired, key: %v, duration: %v seconds\n", key, num)
			time.Sleep(time.Second * time.Duration(num))
			log.Printf("unlock, key: %v\n", key)
		})
	}

	wg.Wait()
	log.Println("finish")
}

func TestEtcdDisLockWithSessionCheck(t *testing.T) {
	// address := "127.0.0.1:2379"
	address := "127.0.0.1:2379, 127.0.0.1:2380, 127.0.0.1:50686"
	key := "lock-3"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "etcd")

	fmt.Println("try to lock")
	lock, err := p9.Lock.Lock(key)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	defer lock.Release()

	recordsToProcess := make([]int, 30)
	for i := range recordsToProcess {
		// 1. DANGER CHECK: Did we lose the lock?
		select {
		case <-lock.Done():
			// Uh oh! The session died. Etcd revoked our lease.
			// We no longer own the lock! We must abort immediately!
			fmt.Printf("CRITICAL: Lost lock halfway through at step %d! Aborting!\n", i)
			assert.Fail(t, "Lost lock halfway through")
			return // Stop processing!
		default:
			// Lock is still healthy, channel is empty. Move on!
		}
		// 2. Perform your heavy work safely
		fmt.Printf("Processing record %d...\n", i)
		time.Sleep(time.Second * 1) // Simulating 1 second of work per record
	}

	fmt.Println("Successfully finished all work while holding the lock!")
}

func TestEtcdDisLockError(t *testing.T) {
	go etcdDisLockErrorLogic(t, "c-1")

	time.Sleep(time.Second * 3)
	go etcdDisLockErrorLogic(t, "c-2")

	time.Sleep(time.Second * 1)
	go etcdDisLockErrorLogic(t, "c-3")

	time.Sleep(time.Second * 1)
	go etcdDisLockErrorLogic(t, "c-4")

	time.Sleep(time.Second * 4)
	go etcdDisLockErrorLogic(t, "c-5")
	go etcdDisLockErrorLogic(t, "c-6")
	go etcdDisLockErrorLogic(t, "c-7")

	time.Sleep(time.Second * 30)
	go etcdDisLockErrorLogic(t, "c-10")
	time.Sleep(time.Second * 10) // wait for c-10
}

func etcdDisLockErrorLogic(t *testing.T, clientName string) {
	// address := "127.0.0.1:2379"
	// address := "127.0.0.1:2379, 127.0.0.1:2380, 127.0.0.1:56209"
	address := "127.0.0.1:2379, 127.0.0.1:2380, 127.0.0.1:50686"
	key := "lock-1"

	p9.Lock.Init(address, time.Second*5, time.Second*3, "etcd")

	fmt.Printf("client: %v, try to lock\n", clientName)
	lock, err := p9.Lock.Lock(key)
	assert.Nil(t, err)
	if err != nil {
		fmt.Printf("cannot get the lock, client: %v\n", clientName)
		return
	}

	defer lock.Release()

	fmt.Printf("client: %v, lock acquired\n", clientName)
	time.Sleep(time.Second * 5)
	fmt.Printf("client: %v, unlock\n", clientName)
}
