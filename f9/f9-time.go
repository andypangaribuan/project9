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

func TimeNow(zone string) time.Time {
	if zone == "" {
		if TimeZone == "" {
			timeNowStr := p9.Conv.Time.ToStr(time.Now(), "yyyy-MM-dd HH:mm:ss.SSSSSS")
			timeNow, _ := p9.Conv.Time.ToTime("yyyy-MM-dd HH:mm:ss.SSSSSS", timeNowStr)
			return timeNow
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

	return time.Now().In(location)
}
