/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package model

import "time"

type RequestHeader struct {
	UID           string `json:"X-UID"`
	Language      string `json:"X-Language"`
	Version       string `json:"X-Version"`
	SvcParent     string `json:"X-SvcParent"`
	RFTime        *time.Time
	RFTimeRaw     string `json:"X-RFTime"`
	HashKey       string `json:"X-HashKey"`
	Authorization string `json:"Authorization"`
	UserHashId    string `json:"X-UserHashId"`
}
