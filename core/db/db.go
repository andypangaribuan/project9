/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type srDb struct{}

func Create() *srDb {
	return &srDb{}
}

type pqInstance struct {
	connRW *srConnection
	connRO *srConnection
}

type pqInstanceTx struct {
	instance   *pqInstance
	tx         *sqlx.Tx
	isCommit   bool
	isRollback bool
	errCommit  error
}

type srConnection struct {
	connStr               string
	driverName            string
	instance              *sqlx.DB
	maxLifeTimeConnection time.Duration
	maxIdleTimeConnection time.Duration
	maxIdleConnection     int
	maxOpenConnection     int
	autoRebind            bool
	unsafeCompatibility   bool
	printSql              bool
}
