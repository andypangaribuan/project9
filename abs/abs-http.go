/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import (
	"time"

	"github.com/andypangaribuan/project9/model"
)

type Http interface {
	Get(url string, header map[string]interface{}, params map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
	Post(url string, header map[string]interface{}, payload map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
	PostData(url string, header map[string]interface{}, data []byte, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
	PostFormData(url string, header map[string]string, files map[string]model.FileData, fields map[string]string) (result []byte, code int, err error)
	Delete(url string, header map[string]interface{}, params map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
}
