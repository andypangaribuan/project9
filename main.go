/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package project9

import (
	"github.com/andypangaribuan/project9/core/check"
	"github.com/andypangaribuan/project9/core/conv"
	"github.com/andypangaribuan/project9/core/crypto"
	"github.com/andypangaribuan/project9/core/db"
	"github.com/andypangaribuan/project9/core/err"
	"github.com/andypangaribuan/project9/core/http"
	"github.com/andypangaribuan/project9/core/json"
	"github.com/andypangaribuan/project9/core/log"
	"github.com/andypangaribuan/project9/core/server"
	"github.com/andypangaribuan/project9/core/util"
	"github.com/andypangaribuan/project9/p9"
)

func Initialize() {
	p9.Init().
		Check(check.Create()).
		Conv(conv.Create()).
		Crypto(crypto.Create()).
		Db(db.Create()).
		Err(err.Create()).
		Http(http.Create()).
		Json(json.Create()).
		Log(log.Create()).
		Server(server.Create()).
		Util(util.Create())
}
