/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "time"

type I9GCS interface {
	Write(destination string, data []byte, timeout ...time.Duration) error
	WriteFromFile(destination string, filePath string, autoIncludeExtension bool, timeout ...time.Duration) error
}
