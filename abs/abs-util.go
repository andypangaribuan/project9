/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import (
	"time"

	"github.com/andypangaribuan/project9/act/actenv"
	"github.com/golang-jwt/jwt"
)

type Util interface {
	GetNanoID(length ...int) string
	GetID25() string
	GetRandom(length int, value string) (string, error)
	BuildJwtToken(privateKey []byte, claims jwt.Claims) (string, error)
	BuildJwtTokenWithPassword(privateKey []byte, password string, claims jwt.Claims) (string, error)
	CreateJwtToken(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte) (string, error)
	CreateJwtTokenWithPassword(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte, password string) (string, error)
	GetJwtClaims(token string, publicKey []byte) (*jwt.StandardClaims, bool, error)
	Base64Encode(data []byte) string
	Base64Decode(value string) ([]byte, error)
}

type UtilEnv interface {
	GetAppEnv(key string) actenv.AppEnv
	GetStr(key string, defValue ...string) string
	GetInt(key string, defValue ...int) int
	GetInt32(key string, defValue ...int32) int32
	GetBool(key string, defValue ...bool) bool
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
