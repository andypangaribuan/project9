/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"github.com/viney-shih/go-lock"
)

type srEnvBase64 struct {
	key  string
	data []byte
}

type srAppEnv struct {
	Value string
}

type srMutex struct {
	mux  lock.Mutex
	name string
}
