/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/act/actenv"
)

func (slf *srEnv) GetAppEnv(key string) actenv.AppEnv {
	return &srAppEnv{
		Value: slf.GetStr(key),
	}
}

func (*srEnv) getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func getEnvDefault[T any](key string, defaultValue ...T) (string, *T) {
	value := strings.TrimSpace(os.Getenv(key))

	switch {
	case value == "" && len(defaultValue) > 0:
		return value, &defaultValue[0]
	case value == "":
		log.Fatalf(`env key "%v" doesn't exists'`, key)
	}

	return value, nil
}

func (*srEnv) GetStr(key string, defaultValue ...string) string {
	strVal, val := getEnvDefault(key, defaultValue...)
	if val != nil {
		return *val
	}

	return strVal
}

func (slf *srEnv) GetInt(key string, defaultValue ...int) int {
	strVal, val := getEnvDefault(key, defaultValue...)
	if val != nil {
		return *val
	}

	value, err := strconv.Atoi(strVal)
	if err != nil {
		log.Fatalf("env key \"%v\" is not int value\nerror:\n%+v", key, err)
	}

	return value
}

func (slf *srEnv) GetInt32(key string, defaultValue ...int32) int32 {
	strVal, val := getEnvDefault(key, defaultValue...)
	if val != nil {
		return *val
	}

	value, err := strconv.ParseInt(strVal, 10, 32)
	if err != nil {
		log.Fatalf(`env key "%v" is not int value\nerror:\n%v`, key, err)
	}

	return int32(value)
}

func (slf *srEnv) GetFloat32(key string, defaultValue ...float32) float32 {
	strVal, val := getEnvDefault(key, defaultValue...)
	if val != nil {
		return *val
	}

	value, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		log.Fatalf("env key \"%v\" is not float32 value\nerror:\n%v", key, err)
	}

	return float32(value)
}

func (slf *srEnv) GetFloat64(key string, defaultValue ...float64) float64 {
	strVal, val := getEnvDefault(key, defaultValue...)
	if val != nil {
		return *val
	}

	value, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		log.Fatalf("env key \"%v\" is not float64 value\nerror:\n%v", key, err)
	}

	return value
}

func (slf *srEnv) GetBool(key string, defaultValue ...bool) bool {
	strVal, val := getEnvDefault(key, defaultValue...)
	if val != nil {
		return *val
	}

	switch strings.ToLower(strVal) {
	case "1", "true":
		return true
	case "0", "false":
		return false
	}

	log.Fatalf("env key \"%v\", from key env key \"%v\" is not a valid boolean value", strVal, key)
	return false
}

func (slf *srEnv) GetBase64(key string) abs.UtilEnvBase64 {
	value := slf.GetStr(key)
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Fatalf("env key \"%v\" is not base64 value\nerror:\n%+v", key, err)
	}

	return &srEnvBase64{
		key:  key,
		data: data,
	}
}

func (slf *srEnvBase64) Key() string {
	return slf.key
}

func (slf *srEnvBase64) Data() []byte {
	return slf.data
}

func (slf *srEnvBase64) String() string {
	return string(slf.data)
}

func (slf *srAppEnv) value() string {
	return strings.ToLower(slf.Value)
}

func (slf *srAppEnv) IsProd() bool {
	val := slf.value()
	return val == "prod" || val == "production"
}

func (slf *srAppEnv) IsStg() bool {
	val := slf.value()
	return val == "stg" || val == "staging"
}

func (slf *srAppEnv) IsDev() bool {
	val := slf.value()
	return val == "dev" || val == "development"
}

func (slf *srAppEnv) IsSandbox() bool {
	val := slf.value()
	return val == "sbx" || val == "sandbox"
}
