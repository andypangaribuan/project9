/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import "github.com/andypangaribuan/project9/abs"

type srGetUpdateHelper struct {
	uh abs.DbUpdateHelper
}

func GetUpdateHelper() abs.DbGetUpdateHelper {
	v := &srGetUpdateHelper{
		uh: NewUpdateHelper(),
	}

	return v
}

func (slf *srGetUpdateHelper) Add(kv map[string]interface{}) abs.DbGetUpdateHelper {
	slf.uh.MapSetAdd(kv)
	return slf
}

func (slf *srGetUpdateHelper) AddIfNotNilOrEmpty(kv map[string]interface{}) abs.DbGetUpdateHelper {
	slf.uh.MapSetAddIfNotNilOrEmpty(kv)
	return slf
}

func (slf *srGetUpdateHelper) Where(kv map[string]interface{}) abs.DbGetUpdateHelper {
	slf.uh.MapWhere(kv)
	return slf
}

func (slf *srGetUpdateHelper) Get() (condition string, pars []interface{}) {
	return slf.uh.Get()
}
