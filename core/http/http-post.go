/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func (slf *srHttp) Post(url string, header map[string]interface{}, payload map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) ([]byte, int, error) {
	payloadData := make([]byte, 0)

	if len(payload) > 0 {
		if v, ok := header["Content-Type"]; ok && v == "application/x-www-form-urlencoded" {
			urlValues := slf.getUrlValues(payload)
			if urlValues != nil {
				urlValues.Encode()
			}
			payloadData = []byte(urlValues.Encode())
		} else {
			data, err := json.Marshal(payload)
			if err != nil {
				return nil, -1, err
			}
			payloadData = data
		}
	}

	return slf.PostData(url, header, payloadData, skipSecurityChecking, timeout)
}

func (slf *srHttp) PostData(url string, header map[string]interface{}, data []byte, skipSecurityChecking bool, timeout *time.Duration) ([]byte, int, error) {
	var payload *bytes.Buffer
	if len(data) > 0 {
		payload = bytes.NewBuffer(data)
		defer payload.Reset()
	}

	request, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, -1, err
	}

	slf.setHeader(request, header)

	return slf.execute(request, skipSecurityChecking, timeout)
}
