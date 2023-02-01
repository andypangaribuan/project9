/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

func (slf *pqInstanceTx) Commit() error {
	return slf.tx.Commit()
}

func (slf *pqInstanceTx) Rollback() error {
	return slf.tx.Rollback()
}
