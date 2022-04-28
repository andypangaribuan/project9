/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package json

import (
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var api jsoniter.API

func init() {
	api = configWithCustomTimeFormat

	setDefaultTimeFormat("2006-01-02 15:04:05.000000", nil)
	addLocaleAlias("-", nil)

	addTimeFormatAlias("date", "2006-01-02")
	addTimeFormatAlias("time", "15:04:05")
	addTimeFormatAlias("full", "2006-01-02 15:04:05")
	addTimeFormatAlias("full-millis", "2006-01-02 15:04:05.000")
	addTimeFormatAlias("full-micros", "2006-01-02 15:04:05.000000")
}

func (*srJson) Marshal(obj interface{}) ([]byte, error) {
	data, err := api.Marshal(obj)
	if err != nil {
		err = errors.WithStack(err)
	}
	return data, err
}

func (*srJson) UnMarshal(data []byte, out interface{}) error {
	err := api.Unmarshal(data, &out)
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (*srJson) Encode(obj interface{}) (string, error) {
	val, err := api.MarshalToString(obj)
	if err != nil {
		err = errors.WithStack(err)
	}
	return val, err
}

func (*srJson) Decode(jsonStr string, out interface{}) error {
	err := api.UnmarshalFromString(jsonStr, &out)
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (*srJson) MapToJson(maps map[string]interface{}) (string, error) {
	val, err := api.MarshalToString(maps)
	if err != nil {
		err = errors.WithStack(err)
	}
	return val, err
}
