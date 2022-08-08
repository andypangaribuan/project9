/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package json

import (
	"github.com/json-iterator/go"
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
	return api.Marshal(obj)
}

func (*srJson) UnMarshal(data []byte, out interface{}) error {
	return api.Unmarshal(data, &out)
}

func (*srJson) Encode(obj interface{}) (string, error) {
	return api.MarshalToString(obj)
}

func (*srJson) Decode(jsonStr string, out interface{}) error {
	return api.UnmarshalFromString(jsonStr, &out)
}

func (*srJson) MapToJson(maps map[string]interface{}) (string, error) {
	return api.MarshalToString(maps)
}
