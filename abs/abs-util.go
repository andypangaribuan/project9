/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

type Util interface {
	GetNanoID(length ...int) (string, error)
}

type UtilEnv interface {
	GetStr(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

type UtilHashId interface {
	Reload()
	Encode(numbers ...int) string
	Add(key, salt string, length int)
	Get(key string) UtilHashId
}
