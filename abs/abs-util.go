/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import (
	"bytes"
	"image"
	"time"

	"github.com/andypangaribuan/project9/act/actenv"
	"github.com/golang-jwt/jwt/v5"
)

type Util interface {
	GetNanoID(length ...int) string
	GetID25() string
	GetXID25() string
	GetXID30() string
	GetXID40() string
	GetRandom(length int, value string) (string, error)
	GetRandomNumber(min, max int) int
	GetRandomAlphabet(length int) string
	GetRandomAlphabetLower(length int) string
	GetRandomAlphabetUpper(length int) string
	GetRandomAlphabetNumber(length int) string
	GetRandomDuration(min int, max int, base time.Duration) time.Duration
	BuildJwtToken(privateKey []byte, claims jwt.Claims) (string, error)
	BuildJwtTokenWithPassword(privateKey []byte, password string, claims jwt.Claims) (string, error)
	CreateJwtToken(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte) (string, error)
	CreateJwtTokenWithPassword(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte, password string) (string, error)
	GetJwtClaims(token string, publicKey []byte) (*jwt.RegisteredClaims, bool, error)
	Base64Encode(data []byte) string
	Base64Decode(value string) ([]byte, error)
	GetExecutionInfo(depth int) (execFunc string, execPath string)
	IsNumberOnly(value string, exclude ...string) bool
	ExtractPhoneNumber(phoneNumber *string) (countryId, countryCode, number string)
	IsEmailValid(email string, verifyDomain ...bool) bool

	ImageBase64Decoder(imageBase64 string) (buff bytes.Buffer, config image.Config, format string, err error)

	// ReflectionSet path: core.util.util-reflection-set.go
	ReflectionSet(obj interface{}, bind map[string]interface{}) error
	ReflectionGet(obj interface{}, fieldName string) (interface{}, error)

	ConcurrentProcess(total, max int, fn func(index int, activeProcess int))
	UniqueConcurrentProcess(total, max int, key func(index int) *string, fn func(index int))
	NewMutex(name string) UtilMutex
	NewNMutex(max int) UtilNMutex
}

type UtilEnv interface {
	GetAppEnv(key string) actenv.AppEnv
	GetStr(key string, defaultValue ...string) string
	GetInt(key string, defaultValue ...int) int
	GetInt32(key string, defaultValue ...int32) int32
	GetInt64(key string, defaultValue ...int64) int64
	GetFloat32(key string, defaultValue ...float32) float32
	GetFloat64(key string, defaultValue ...float64) float64
	GetBool(key string, defaultValue ...bool) bool
	GetBase64(key string) UtilEnvBase64
	GetStrArr(key string, separator string, defaultValue ...[]string) []string
	GetMapStrStr(key string, keyValSeparator string, itemSeparator string, defaultValue ...map[string]string) map[string]string
}

type UtilEnvBase64 interface {
	Key() string
	Data() []byte
}

type UtilHashId interface {
	Reload()
	Encode(numbers ...int) string
	Decode(hashed string) ([]int, error)
	Encode64(numbers ...int64) string
	DecodeInt64(hashed string) ([]int64, error)
	Add(key, salt string, length int)
	Get(key string) UtilHashId
}

type UtilMutex interface {
	Sleep(duration ...time.Duration)
	Lock(timeout ...time.Duration) (isTimeout bool)
	Unlock()
	Exec(timeout *time.Duration, fn func()) (executed bool, panicErr error)
	FExec(timeoutLock *time.Duration, timeoutFunc time.Duration, fn func()) (executed bool, isTimeout bool, panicErr error)
	Func(timeout time.Duration, fn func()) (isTimeout bool, panicErr error)
}

type UtilNMutex interface {
	Lock(key string, totalTry ...int) (locked bool)
	Unlock(key string)
}
