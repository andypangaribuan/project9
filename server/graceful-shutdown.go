/*
 * Copyright (c) 2026.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andypangaribuan/project9/p9"
)

var (
	gsMtx                  sync.Mutex
	isGracefulShutdownImpl bool
)

func gracefulShutdown() {
	gsMtx.Lock()
	defer gsMtx.Unlock()

	if isGracefulShutdownImpl {
		return
	}

	isGracefulShutdownImpl = true
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		// Block until we receive a signal
		<-c

		time.Sleep(time.Second * 10)
		p9.Lock.Close()
		os.Exit(0) // <-- 0 indicates a clean intentional exit
	}()
}
