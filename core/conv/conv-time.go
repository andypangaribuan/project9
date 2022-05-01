/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	layoutTimeDate = "2006-01-02"
	layoutTimeFull = "2006-01-02 15:04:05"
)

var replacer = [][]string{
	{"yyyy", "2006"},
	{"MM", "01"},
	{"dd", "02"},
	{"HH", "15"},
	{"mm", "04"},
	{"ss", "05"},
	{"SSSSSS", "000000"},
	{"SSSSS", "00000"},
	{"SSSS", "0000"},
	{"SSS", "000"},
	{"SS", "00"},
	{"S", "0"},
}

func (*srTime) ToStr(tm time.Time, format string) string {
	for _, arr := range replacer {
		format = strings.Replace(format, arr[0], arr[1], -1)
	}
	return tm.Format(format)
}

func (*srTime) ToStrDate(tm time.Time) string {
	return tm.Format(layoutTimeDate)
}

func (*srTime) ToStrFull(tm time.Time) string {
	return tm.Format(layoutTimeFull)
}

func (*srTime) ToTime(layout string, value string) (tm time.Time, err error) {
	for _, arr := range replacer {
		layout = strings.Replace(layout, arr[0], arr[1], -1)
	}

	tm, err = time.Parse(layout, value)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
