/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

type UtilEnv interface {
	GetStr(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

type UtilHashId interface {
	Reload()
	Encode(numbers []int) string
	EncodeId(number int) string
}
