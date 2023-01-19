/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package fc

import (
	"database/sql/driver"
)

// Scan implements the sql.Scanner interface for database deserialization.
func (slf *FCT) Scan(value interface{}) error {
	err := slf.vd.Scan(value)
	if err != nil {
		return err
	}

	slf.set(slf.vd)
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (slf FCT) Value() (driver.Value, error) {
	return slf.vd.String(), nil
}
