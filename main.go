/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package project9

import (
	"github.com/andypangaribuan/project9/core/check"
	"github.com/andypangaribuan/project9/core/conv"
	"github.com/andypangaribuan/project9/core/db"
	"github.com/andypangaribuan/project9/core/json"
	"github.com/andypangaribuan/project9/core/log"
	"github.com/andypangaribuan/project9/core/util"
	"github.com/andypangaribuan/project9/p9"
)

//goland:noinspection ALL
func Initialize() {
	p9.Init().
		Check(check.Create()).
		Conv(conv.Create()).
		Db(db.Create()).
		Json(json.Create()).
		Log(log.Create()).
		Util(util.Create())
}
