/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Util interface {
	GetNanoID(length ...int) (string, error)
	CreateJwtToken(subject, id string, expiresAt, issuedAt, notBefore time.Time, pemKey []byte) (string, error)
	GetJwtClaims(token string, publicKey []byte) (*jwt.StandardClaims, bool, error)
}

type UtilEnv interface {
	GetStr(key string) string
	GetInt(key string) int
	GetInt32(key string) int32
	GetBool(key string) bool
	GetBase64(key string) UtilEnvBase64
}

type UtilEnvBase64 interface {
	Key() string
	Data() []byte
}

type UtilHashId interface {
	Reload()
	Encode(numbers ...int) string
	Encode64(numbers ...int64) string
	Add(key, salt string, length int)
	Get(key string) UtilHashId
}
