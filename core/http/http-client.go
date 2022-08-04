/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/andypangaribuan/project9/p9"
	"github.com/gojektech/heimdall/httpclient"
)

func (slf *srHttpClient) Do(request *http.Request) (*http.Response, error) {
	return slf.client.Do(request)
}

func getHttpClient(skipSecurityChecking bool, timeout *time.Duration) (client *httpclient.Client) {
	httpTimeout := 5 * time.Minute

	if p9.Conf.DefaultHttpRequestTimeout != nil {
		httpTimeout = *p9.Conf.DefaultHttpRequestTimeout
	}
	if timeout != nil {
		httpTimeout = *timeout
	}

	if !skipSecurityChecking {
		client = httpclient.NewClient(httpclient.WithHTTPTimeout(httpTimeout))
	} else {
		client = httpclient.NewClient(
			httpclient.WithHTTPClient(&srHttpClient{
				client: http.Client{
					Timeout: httpTimeout,
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				},
			}))
	}

	return
}
