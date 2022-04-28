/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/jmoiron/sqlx"
)

type srDb struct {
	abs.Db
}

func Create() *srDb {
	return &srDb{}
}

type postgresInstance struct {
	abs.DbPostgresInstance
	data *connModel
}

type connModel struct {
	connStr               string
	driverName            string
	instance              *sqlx.DB
	maxLifeTimeConnection time.Duration
	maxIdleConnection     int
	maxOpenConnection     int
	autoRebind            bool
	unsafeCompatibility   bool
}
