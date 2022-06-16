/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type srInit struct{}

//goland:noinspection ALL
var (
	Check *srCheck
	Conf  *srConf
	Conv  *srConv
	Db    *srDb
	Json  *srJson
	Log   *srLog
	Util  *srUtil
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

func (slf *srInit) Check(fnStr abs.CheckStr) *srInit {
	Check = &srCheck{
		Str: &srCheckStr{fnStr},
	}
	return slf
}

func (slf *srInit) Conv(fnTime abs.ConvTime) *srInit {
	Conv = &srConv{
		Time: &srConvTime{fnTime},
	}
	return slf
}

func (slf *srInit) Db(fnDb abs.Db) *srInit {
	Db = &srDb{fnDb}
	return slf
}

func (slf *srInit) Json(fnJson abs.Json) *srInit {
	Json = &srJson{fnJson}
	return slf
}

func (slf *srInit) Log(fnLog abs.Log) *srInit {
	Log = &srLog{fnLog}
	return slf
}

func (slf *srInit) Util(fnUtil abs.Util, fnEnv abs.UtilEnv, fnHashId abs.UtilHashId) *srInit {
	Util = &srUtil{
		absUtil: fnUtil,
		Env:     &srUtilEnv{fnEnv},
		HashId:  &srUtilHashId{fnHashId},
	}
	return slf
}
