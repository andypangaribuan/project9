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
	type model struct {
		id  string
		key string
	}

	var (
		trx = p9.Util.NewNMutex(4)
		ls  = []model{
			{id: "a", key: "1"},
			{id: "b", key: "2"},
			{id: "c", key: "3"},
			{id: "d", key: "1"},
			{id: "e", key: "5"},
		}
		total = len(ls)
		max   = 4
		key   = func(index int) *string {
			m := ls[index]
			locked := trx.Lock(m.key, 3)
			return f9.Ternary(locked, &m.key, nil)
		}
		fn = func(index int) {
			m := ls[index]
			log.Println("start:", m.id, m.key)
			defer func() {
				trx.Unlock(m.key)
				log.Println("done:", m.id, m.key)
			}()

			time.Sleep(time.Second * 3)
		}
	)

	p9.Util.UniqueConcurrentProcess(total, max, key, fn)
}
