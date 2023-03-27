/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "time"

type srConf struct {
	NanoIdLength int
	HashIdSalt   string
	HashIdLength int
	K8sEnvName   string

	DefaultHttpRequestTimeout *time.Duration
}
