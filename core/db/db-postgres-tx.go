/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

func (slf *pqInstanceTx) Commit() error {
	if slf.isCommit || slf.isRollback || slf.tx == nil {
		return nil
	}

	slf.isCommit = true
	return slf.tx.Commit()
}

func (slf *pqInstanceTx) Rollback() error {
	if slf.isCommit || slf.isRollback || slf.tx == nil {
		return nil
	}

	slf.isRollback = true
	return slf.tx.Rollback()
}
