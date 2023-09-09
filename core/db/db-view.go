/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
	"github.com/andypangaribuan/project9/p9"
	"github.com/lib/pq"
)

type View struct {
	dbInstance     abs.DbInstance
	printUnsafeErr bool
}

func (slf *View) New(db abs.DbInstance, printUnsafeErr bool) *View {
	slf.dbInstance = db
	slf.printUnsafeErr = printUnsafeErr
	return slf
}

func (slf *View) onUnsafe(unsafe *model.DbUnsafeSelectError) {
	if unsafe != nil && slf.printUnsafeErr {
		log.Printf("[%v] db.unsafe.select.error:\nerror: %v\nsql-query: %v\nsql-pars: %v\ntrace: %v\n",
			p9.Conv.Time.ToStrFull(f9.TimeNow()),
			f9.TernaryFnB(unsafe.LogMessage == nil, "nil", func() string { return *unsafe.LogMessage }),
			unsafe.SqlQuery,
			unsafe.SqlPars,
			f9.TernaryFnB(unsafe.LogTrace == nil, "nil", func() string { return *unsafe.LogTrace }),
		)
	}
}

func (slf *View) Select(out interface{}, sqlQuery string, sqlPars ...interface{}) error {
	unsafe, _, err := slf.dbInstance.Select(true, out, sqlQuery, sqlPars...)
	slf.onUnsafe(unsafe)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			msg := strings.TrimSpace(e.Message)
			if msg != "" {
				msg += "\n"
			}
			msg += fmt.Sprintf("sql: %v\npars: %v", sqlQuery, sqlPars)
			e.Message = msg
			err = e
		}
		return err
	}

	return nil
}
