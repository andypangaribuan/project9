/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type absLock interface {
	abs.Lock
}

type srLock struct {
	absLock
}
