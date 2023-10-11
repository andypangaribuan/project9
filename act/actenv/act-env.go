/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package actenv

type AppEnv interface {
	IsProd() bool
	IsPreProd() bool
	IsStg() bool
	IsDev() bool
	IsSandbox() bool
	IsPPS() bool
}
