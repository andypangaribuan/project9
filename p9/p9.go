/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type srInit struct{}

//goland:noinspection ALL
var (
	Conf *srConf
	Conv *srConv
	Db   *srDb
	Json *srJson
	Util *srUtil
)

func init() {
	Conf = &srConf{
		NanoIdLength: 60,
		HashIdSalt:   "Project9",
		HashIdLength: 60,
	}
}

func Init() *srInit {
	return nil
}

func (slf *srInit) Conv(fnTime abs.ConvTime) *srInit {
	Conv = &srConv{
		Time: &srConvTime{ConvTime: fnTime},
	}
	return slf
}

func (slf *srInit) Db(fnDb abs.Db) *srInit {
	Db = &srDb{Db: fnDb}
	return slf
}

func (slf *srInit) Json(fnJson abs.Json) *srInit {
	Json = &srJson{Json: fnJson}
	return slf
}

func (slf *srInit) Util(fnUtil abs.Util, fnEnv abs.UtilEnv, fnHashId abs.UtilHashId) *srInit {
	Util = &srUtil{
		Util:   fnUtil,
		Env:    &srUtilEnv{UtilEnv: fnEnv},
		HashId: &srUtilHashId{UtilHashId: fnHashId},
	}
	return slf
}
