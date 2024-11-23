/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

type FetchOpt struct {
	EndQuery               string
	WithoutDeletedAtIsNull bool
	ForceRW                bool
}

type srFetchOpt struct {
	rwForce bool
}

type Update struct {
	Set       string
	Where     string
	SetPars   []interface{}
	WherePars []interface{}
}
