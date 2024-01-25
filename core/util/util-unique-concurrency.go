/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"sync"
	"time"
)

type srUniqueConcurrency struct {
	sr     *srUtil
	mx     sync.Mutex
	max    int
	active map[string]interface{}
	queue  []int
	key    func(index int) *string
	fn     func(index int)
}

func (slf *srUtil) UniqueConcurrentProcess(total, max int, key func(index int) *string, fn func(index int)) {
	c := &srUniqueConcurrency{
		sr:     slf,
		active: make(map[string]interface{}, 0),
		max:    max,
		key:    key,
		fn:     fn,
	}

	c.queue = make([]int, 0)
	for i := 0; i < total; i++ {
		c.queue = append(c.queue, i)
	}

	c.execute()
}

func (slf *srUniqueConcurrency) sleep() {
	time.Sleep(slf.sr.GetRandomDuration(1, 10, time.Millisecond))
}

func (slf *srUniqueConcurrency) getTotalQueue() int {
	slf.mx.Lock()
	defer slf.mx.Unlock()
	return len(slf.queue)
}

func (slf *srUniqueConcurrency) getTotalActive() int {
	slf.mx.Lock()
	defer slf.mx.Unlock()
	return len(slf.active)
}

func (slf *srUniqueConcurrency) removeQueueAt(index int) {
	slf.queue = append(slf.queue[:index], slf.queue[index+1:]...)
}

func (slf *srUniqueConcurrency) getNext() (index int, key string) {
	slf.mx.Lock()
	defer slf.mx.Unlock()

	for qi, index := range slf.queue {
		_key := slf.key(index)
		if _key == nil {
			continue
		}

		if _, ok := slf.active[*_key]; !ok {
			slf.active[*_key] = nil
			slf.removeQueueAt(qi)
			return index, *_key
		}
	}

	return -1, ""
}

func (slf *srUniqueConcurrency) execute() {
	for {
		totalQueue := slf.getTotalQueue()
		totalActive := slf.getTotalActive()

		if totalQueue == 0 && totalActive == 0 {
			break
		}

		if totalQueue == 0 || totalActive == slf.max {
			slf.sleep()
			continue
		}

		nextIndex, key := slf.getNext()
		if nextIndex == -1 {
			slf.sleep()
			continue
		}

		go slf.call(nextIndex, key)
	}
}

func (slf *srUniqueConcurrency) call(index int, key string) {
	slf.fn(index)
	slf.mx.Lock()
	defer slf.mx.Unlock()
	delete(slf.active, key)
}
