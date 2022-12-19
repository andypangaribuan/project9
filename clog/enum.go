/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package clog

type Severity int

const (
	Info Severity = iota + 1
	Warning
	Error
)

func (s Severity) String() string {
	return [...]string{"info", "warning", "error"}[s-1]
}

func (s Severity) Index() int {
	return int(s)
}
