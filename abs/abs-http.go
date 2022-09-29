/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

import "time"

type Http interface {
	Get(url string, header map[string]interface{}, payload map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
	Post(url string, header map[string]interface{}, payload map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
	PostData(url string, header map[string]interface{}, data []byte, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error)
	Delete(url string, header map[string]interface{}, params map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) ([]byte, int, error)
}
