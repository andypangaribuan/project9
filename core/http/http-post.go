/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/model"
)

func (slf *srHttp) Post(url string, header map[string]interface{}, payload map[string]interface{}, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error) {
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

func (slf *srHttp) PostData(url string, header map[string]interface{}, data []byte, skipSecurityChecking bool, timeout *time.Duration) (result []byte, code int, err error) {
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

func (slf *srHttp) PostFormData(url string, header map[string]string, files map[string]model.FileData, fields map[string]string) (result []byte, code int, err error) {
	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)
	buffLs := make([]*bytes.Reader, 0)

	defer func() {
		zeroData := make([]byte, 0)
		for _, buff := range buffLs {
			if buff != nil {
				buff.Reset(zeroData)
			}
		}
	}()

	// add files
	for fieldName, fileData := range files {
		buff := bytes.NewReader(fileData.Data)
		buffLs = append(buffLs, buff)

		partWriter, err := bodyWriter.CreateFormFile(fieldName, fileData.Name)
		if err != nil {
			return nil, -1, err
		}

		_, err = io.Copy(partWriter, buff)
		if err != nil {
			return nil, -1, err
		}
	}

	// add field value
	for fieldName, value := range fields {
		err := bodyWriter.WriteField(fieldName, value)
		if err != nil {
			return nil, -1, err
		}
	}

	err = bodyWriter.Close()
	if err != nil {
		return nil, -1, err
	}

	// header checking
	haveHeaderAccept := false
	for key := range header {
		if strings.ToLower(key) == "accept" {
			haveHeaderAccept = true
		}
	}

	if !haveHeaderAccept {
		header["Accept"] = "application/json"
	}
	header["Content-Type"] = bodyWriter.FormDataContentType()

	return doHttpRequest("POST", url, header, body)
}

func doHttpRequest(method string, url string, header map[string]string, body io.Reader) (result []byte, code int, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, -1, err
	}

	for key, value := range header {
		request.Header.Add(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, -1, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	return resData, res.StatusCode, err
}
