/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"context"
	"errors"
	"time"
)

func Retry[T any](executor func() (T, error), retries int, delay time.Duration, ctx ...context.Context) (T, error) {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	if retries < 1 {
		for {
			response, err := executor()
			if err == nil {
				return response, err
			}

			select {
			case <-time.After(delay):
			case <-c.Done():
				return *new(T), errors.New("force-close")
			}
		}
	}

	for r := 1; ; r++ {
		response, err := executor()
		if err == nil || r >= retries {
			return response, err
		}

		select {
		case <-time.After(delay):
		case <-c.Done():
			return *new(T), errors.New("force-close")
		}
	}
}
