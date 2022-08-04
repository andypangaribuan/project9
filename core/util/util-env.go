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
	"github.com/pkg/errors"
)

func (slf *srEnv) GetAppEnv(key string) actenv.AppEnv {
	return &srAppEnv{
		Value: slf.GetStr(key),
	}
}

func (*srEnv) GetStr(key string) string {
	value := os.Getenv(key)
	value = strings.TrimSpace(value)
	if value == "" {
		log.Fatalf("env key \"%v\" not found", key)
	}
	return value
}

func (slf *srEnv) GetInt(key string) int {
	value, err := strconv.Atoi(slf.GetStr(key))
	if err != nil {
		err = errors.WithStack(err)
		log.Fatalf("env key \"%v\" is not int value\nerror:\n%+v", key, err)
	}
	return value
}

func (slf *srEnv) GetInt32(key string) int32 {
	value, err := strconv.ParseInt(slf.GetStr(key), 10, 32)
	if err != nil {
		err = errors.WithStack(err)
		log.Fatalf("env key \"%v\" is not int value\nerror:\n%+v", key, err)
	}
	return int32(value)
}

func (slf *srEnv) GetBool(key string) bool {
	value := strings.ToLower(slf.GetStr(key))
	if value == "1" || value == "true" {
		return true
	}
	if value == "0" || value == "false" {
		return false
	}
	log.Fatalf("env key \"%v\", from key env key \"%v\" is not a valid boolean value", value, key)
	return false
}

func (slf *srEnv) GetBase64(key string) abs.UtilEnvBase64 {
	value := slf.GetStr(key)
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		err = errors.WithStack(err)
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
