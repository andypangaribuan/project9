/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

func (*srHttp) getInterfaceString(value interface{}) string {
	res := ""

	if value != nil {
		vType := fmt.Sprintf("%T", value)
		if vType[:1] == "*" {
			refVal := reflect.ValueOf(value)
			if !refVal.IsNil() {
				res = fmt.Sprintf("%v", refVal.Elem())
			}
		} else {
			res = fmt.Sprintf("%v", value)
		}
	}

	return res
}

func (slf *srHttp) setHeader(request *http.Request, header map[string]interface{}) {
	isContentTypeAdded := false
	if header != nil {
		for key, val := range header {
			request.Header.Add(key, slf.getInterfaceString(val))
			if key == "Content-Type" {
				isContentTypeAdded = true
			}
		}
	}
	if !isContentTypeAdded {
		request.Header.Set("Content-Type", "application/json")
	}
}

func (slf *srHttp) setParams(request *http.Request, params map[string]interface{}) {
	if params != nil {
		query := request.URL.Query()
		for key, val := range params {
			query.Add(key, slf.getInterfaceString(val))
		}
		request.URL.RawQuery = query.Encode()
	}
}

func (slf *srHttp) getUrlValues(payload map[string]interface{}) *url.Values {
	if len(payload) == 0 {
		return nil
	}

	data := &url.Values{}
	for key, val := range payload {
		data.Set(key, slf.getInterfaceString(val))
	}
	return data
}

func (slf *srHttp) execute(request *http.Request, skipSecurityChecking bool, timeout *time.Duration) ([]byte, int, error) {
	client := getHttpClient(skipSecurityChecking, timeout)
	response, err := client.Do(request)
	if err != nil {
		return nil, -1, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, -1, err
	}

	return responseBody, response.StatusCode, nil
}
