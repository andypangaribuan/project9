/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import (
	"net/http"
	"time"
)

func (slf *srHttp) Get(url string, header map[string]interface{}, params map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, -1, err
	}

	slf.setHeader(request, header)
	slf.setParams(request, params)

	return slf.execute(request, skipSecurityChecking, timeout)
}
