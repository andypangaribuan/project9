/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"time"

	"github.com/andypangaribuan/project9/p9"
)

var TimeZone string
var timeZones map[string]*time.Location

func init() {
	timeZones = make(map[string]*time.Location, 0)
}

func TimeNow(timezone ...string) time.Time {
	location := getTimeLocation(timezone...)
	return time.Now().In(location)
}

func TimeRemoveTimeZone(tm time.Time) time.Time {
	millis := p9.Conv.Time.ToStrMillis(tm)
	timeWithoutTimeZone, _ := p9.Conv.Time.ToTimeMillis(millis)
	return timeWithoutTimeZone
}

func CustomTime(year, month, day, hour, min, sec, nsec int, timezone ...string) time.Time {
	location := getTimeLocation(timezone...)
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, location)
}

func getTimeLocation(timezone ...string) *time.Location {
	zone := ""
	if len(timezone) > 0 && timezone[0] != "" {
		zone = timezone[0]
	}

	if zone == "" {
		if TimeZone == "" {
			return time.UTC
		}
		zone = TimeZone
	}

	var location *time.Location
	if loc, ok := timeZones[zone]; !ok {
		_loc, _ := time.LoadLocation(zone)
		location = _loc
		timeZones[zone] = _loc
	} else {
		location = loc
	}

	return location
}
