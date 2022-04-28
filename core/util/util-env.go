/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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
