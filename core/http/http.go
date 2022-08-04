/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package http

import "net/http"

type srHttp struct{}

type srHttpClient struct {
	client http.Client
}

func Create() *srHttp {
	return &srHttp{}
}
