/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package test

import (
	"log"
	"testing"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

func TestUniqueConcurrentProcess(t *testing.T) {
	var (
		trx   = p9.Util.NewNMutex(4)
		ls    = []string{"1", "2", "3", "1", "4"}
		total = len(ls)
		max   = 4
		key   = func(index int) *string {
			id := ls[index]
			locked := trx.Lock(id, 3)
			return f9.Ternary(locked, &id, nil)
		}
		fn = func(index int) {
			id := ls[index]
			log.Println("start:", id)
			defer func ()  {
				trx.Unlock(id)
				log.Println("done:", id)
			}()

			time.Sleep(time.Second * 3)
		}
	)

	p9.Util.UniqueConcurrentProcess(total, max, key, fn)
}
