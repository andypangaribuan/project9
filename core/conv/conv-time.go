/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

import (
	"fmt"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/f9"
)

const (
	layoutTimeDate   = "2006-01-02"
	layoutTimeFull   = "2006-01-02 15:04:05"
	layoutTimeMillis = "2006-01-02 15:04:05.000"
	layoutTimeMicro  = "2006-01-02 15:04:05.000000"
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

func (*srTime) ToStrMillis(tm time.Time) string {
	return tm.Format(layoutTimeMillis)
}

func (*srTime) ToStrMicro(tm time.Time) string {
	return tm.Format(layoutTimeMicro)
}

func (*srTime) ToStrRFC3339(tm time.Time) string {
	return tm.Format(time.RFC3339)
}

func (*srTime) ToStrRFC3339MilliSecond(tm time.Time) string {
	return tm.Format("2006-01-02T15:04:05.000Z07:00")
}

func (*srTime) ToStrRFC3339MicroSecond(tm time.Time) string {
	return tm.Format("2006-01-02T15:04:05.000000Z07:00")
}

func (*srTime) ToTime(layout string, value string) (time.Time, error) {
	for _, arr := range replacer {
		layout = strings.Replace(layout, arr[0], arr[1], -1)
	}
	return time.Parse(layout, value)
}

func (slf *srTime) ToTimeDate(value string) (time.Time, error) {
	return time.Parse(layoutTimeDate, value)
}

func (slf *srTime) ToTimeFull(value string) (time.Time, error) {
	return time.Parse(layoutTimeFull, value)
}

func (slf *srTime) ToTimeMillis(value string) (time.Time, error) {
	return time.Parse(layoutTimeMillis, value)
}

func (slf *srTime) ToTimeMicro(value string) (time.Time, error) {
	return time.Parse(layoutTimeMicro, value)
}

func (slf *srTime) ToTimeRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func (slf *srTime) ToTimeRFC3339MilliSecond(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z07:00", value)
}

func (slf *srTime) ToTimeRFC3339MicroSecond(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000000Z07:00", value)
}

func (slf *srTime) GetTimeZone(tm time.Time) string {
	hour := 3600 // second
	minute := 60 // second
	isPlus := true
	_, offset := tm.Zone()

	if offset < 0 {
		isPlus = false
		offset *= -1
	}

	h := offset / hour
	offset = offset - (hour * h)

	m := offset / minute

	hs := fmt.Sprintf("%v", h)
	if len(hs) == 1 {
		hs = "0" + hs
	}

	ms := fmt.Sprintf("%v", m)
	if len(ms) == 1 {
		ms = "0" + ms
	}

	zone := f9.Ternary(isPlus, "+", "-")
	zone += hs + ":" + ms

	return zone
}

func (slf *srTime) RemovePart(tm time.Time, part ...string) time.Time {
	for _, p := range part {
		switch p {
		case "yyyy":
			tm = time.Date(0, tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(), tm.Location())
		case "MM":
			tm = time.Date(tm.Year(), 1, tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(), tm.Location())
		case "dd":
			tm = time.Date(tm.Year(), tm.Month(), 1, tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(), tm.Location())
		case "HH":
			tm = time.Date(tm.Year(), tm.Month(), tm.Day(), 0, tm.Minute(), tm.Second(), tm.Nanosecond(), tm.Location())
		case "mm":
			tm = time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), 0, tm.Second(), tm.Nanosecond(), tm.Location())
		case "ss":
			tm = time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), 0, tm.Nanosecond(), tm.Location())
		case "ns":
			tm = time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, tm.Location())
		}
	}

	return tm
}
