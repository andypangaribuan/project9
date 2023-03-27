/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import (
	"time"

	"github.com/andypangaribuan/project9/abs"
)

type srInit struct{}

//goland:noinspection ALL
var (
	Check  *srCheck
	Conf   *srConf
	Conv   *srConv
	Crypto *srCrypto
	Db     *srDb
	Err    *srErr
	Http   *srHttp
	Json   *srJson
	Log    *srLog
	Server *srServer
	Util   *srUtil
)

func init() {
	defaultHttpRequestTimeout := 2 * time.Minute
	Conf = &srConf{
		NanoIdLength:              60,
		HashIdSalt:                "Project9",
		HashIdLength:              60,
		K8sEnvName:                "ZX_ENV",
		DefaultHttpRequestTimeout: &defaultHttpRequestTimeout,
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

func (slf *srInit) Conv(fnConv abs.Conv, fnTime abs.ConvTime, fnProto abs.ConvProto) *srInit {
	Conv = &srConv{
		absConv: fnConv,
		Time:    &srConvTime{fnTime},
		Proto:   &srConvProto{fnProto},
	}
	return slf
}

func (slf *srInit) Crypto(fnAES abs.CryptoAES, fnMD5 abs.CryptoMD5, fnSHA256 abs.CryptoSHA256) *srInit {
	Crypto = &srCrypto{
		AES:    &srCryptoAES{fnAES},
		MD5:    &srCryptoMD5{fnMD5},
		SHA256: &srCryptoSHA256{fnSHA256},
	}
	return slf
}

func (slf *srInit) Db(fnDb abs.Db) *srInit {
	Db = &srDb{fnDb}
	return slf
}

func (slf *srInit) Err(fnErr abs.Err) *srInit {
	Err = &srErr{fnErr}
	return slf
}

func (slf *srInit) Http(fnHttp abs.Http) *srInit {
	Http = &srHttp{fnHttp}
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

func (slf *srInit) Server(fn abs.Server) *srInit {
	Server = &srServer{fn}
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
