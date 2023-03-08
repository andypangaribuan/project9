/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"sync"
	"time"

	"github.com/andypangaribuan/project9/f9"
)

type srConcurrency struct {
	mx     sync.Mutex
	max    int
	total  int
	active int
	fn     func(index int)
}

func (*srUtil) ConcurrentProcess(total, max int, fn func(index int)) {
	c := &srConcurrency{
		active: 0,
		total:  total,
		max:    max,
		fn:     fn,
	}

	c.start()
}

func (slf *srConcurrency) start() {
	n := 0
	for i := 0; i < slf.total; i++ {
		if slf.active >= slf.max {
			for {
				time.Sleep(time.Millisecond * 10)
				if slf.active < slf.max {
					break
				}
			}
		}

		n++
		slf.addActive(1)
		idx := f9.Ptr(i)
		go slf.execute(*idx)
	}

	for {
		time.Sleep(time.Millisecond * 10)
		if slf.active == 0 {
			println("n: ", n)
			break
		}
	}
}

func (slf *srConcurrency) execute(index int) {
	slf.fn(index)
	slf.addActive(-1)
}

func (slf *srConcurrency) addActive(add int) {
	slf.mx.Lock()
	defer slf.mx.Unlock()
	slf.active += add
}
