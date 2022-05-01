/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type absUtil interface {
	abs.Util
}

type absUtilEnv interface {
	abs.UtilEnv
}

type absUtilHashId interface {
	abs.UtilHashId
}

type srUtil struct {
	absUtil
	Env    *srUtilEnv
	HashId *srUtilHashId
}

type srUtilEnv struct {
	absUtilEnv
}

type srUtilHashId struct {
	absUtilHashId
}
