/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/andypangaribuan/project9"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

func TestRetry(t *testing.T) {
	var i int = 0

	f := func() (*string, error) {
		i++
		if i == 10 {
			v := "value"
			return &v, nil
		}

		log.Printf("invalid i: %v\n", i)
		return nil, fmt.Errorf("invalid i : %v", i)
	}

	r, err := f9.Retry(f, 4, time.Second)
	if err != nil {
		log.Printf("error: %v\n", err)
		t.Error(err)
		return
	}

	log.Printf("done: %v\n", *r)
	log.Printf("%v", *r)
}

func TestXID(t *testing.T) {
	project9.Initialize()
	// tm := f9.TimeNow()

	// t.Log(tm.Unix())
	// t.Log(tm.UnixMilli())
	// t.Log(tm.UnixMicro())
	// t.Log(tm.UnixNano())

	hex := fmt.Sprintf("%v", f9.TimeNow().UnixMicro())
	nine := p9.Util.GetRandomAlphabetNumber(9)
	xid := hex + nine
	log.Println(xid)
}
