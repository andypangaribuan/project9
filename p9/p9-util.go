/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type srUtil struct {
	abs.Util
	Env    *srUtilEnv
	HashId *srUtilHashId
}

type srUtilEnv struct {
	abs.UtilEnv
}

type srUtilHashId struct {
	abs.UtilHashId
}
