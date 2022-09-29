/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import (
	"net/http"
	"time"
)

func (slf *srHttp) Delete(url string, header map[string]interface{}, params map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) ([]byte, int, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, -1, err
	}

	slf.setHeader(request, header)
	slf.setParams(request, params)

	return slf.execute(request, skipSecurityChecking, timeout)
}
