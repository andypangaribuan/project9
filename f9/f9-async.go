/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import "context"

type Future interface {
	Await() interface{}
}

type future struct {
	await func(ctx context.Context) interface{}
}

type srAsync struct {
	result interface{}
	c      chan struct{}
}

func (slf *srAsync) future() Future {
	return future{
		await: func(ctx context.Context) interface{} {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-slf.c:
				return slf.result
			}
		},
	}
}

func (f future) Await() interface{} {
	return f.await(context.Background())
}

func Async(fn func() interface{}) Future {
	sr := &srAsync{
		c: make(chan struct{}),
	}
	go func() {
		defer close(sr.c)
		sr.result = fn()
	}()
	return sr.future()
}

func Async1[T1 any](fn func(p1 T1) interface{}, p1 T1) Future {
	sr := &srAsync{
		c: make(chan struct{}),
	}
	go func() {
		defer close(sr.c)
		sr.result = fn(p1)
	}()
	return sr.future()
}

func Async2[T1 any, T2 any](fn func(p1 T1, p2 T2) interface{}, p1 T1, p2 T2) Future {
	sr := &srAsync{
		c: make(chan struct{}),
	}
	go func() {
		defer close(sr.c)
		sr.result = fn(p1, p2)
	}()
	return sr.future()
}

func Async3[T1 any, T2 any, T3 any](fn func(p1 T1, p2 T2, p3 T3) interface{}, p1 T1, p2 T2, p3 T3) Future {
	sr := &srAsync{
		c: make(chan struct{}),
	}
	go func() {
		defer close(sr.c)
		sr.result = fn(p1, p2, p3)
	}()
	return sr.future()
}
