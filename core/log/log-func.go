/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package log

import (
	"fmt"
	"runtime"
)

func (*srLog) GetSourceCall() string {
	pc, filename, line, _ := runtime.Caller(1)
	return fmt.Sprintf("[source:%v:%v:%v]", runtime.FuncForPC(pc).Name(), filename, line)
}
