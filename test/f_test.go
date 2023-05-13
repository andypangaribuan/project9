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

	"github.com/andypangaribuan/project9/f9"
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
