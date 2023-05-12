/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import "sync"

type srEnvBase64 struct {
	key  string
	data []byte
}

type srAppEnv struct {
	Value string
}

type srMutex struct {
	mtx *sync.Mutex
}
