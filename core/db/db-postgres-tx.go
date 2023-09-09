/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
)

const min int64 = 10
const max int64 = 300

func (slf *pqInstanceTx) Commit() (err error) {
	if slf == nil {
		return
	}

	if slf.isCommit || slf.isRollback || slf.tx == nil {
		return nil
	}

	err = slf.tx.Commit()
	slf.isCommit = err == nil
	return


	// iteration := 0
	// const maxTry = 3

	// for {
	// 	iteration++
	// 	err = slf.tx.Commit()

	// 	if err == nil {
	// 		break
	// 	}

	// 	if err == sql.ErrConnDone {
	// 		break
	// 	}

	// 	if err == sql.ErrTxDone {
	// 		err = nil
	// 		break
	// 	}

	// 	if iteration >= maxTry {
	// 		break
	// 	}

	// 	time.Sleep(time.Microsecond * time.Duration(rand.Int63n(max-min)+min))
	// }

	// slf.isCommit = err == nil
	// return

	// slf.isCommit = true
	// err = slf.tx.Commit()
	// if err != nil {
	// 	log.Printf("db.tx.commit error: %+v\n", err)
	// }

	// slf.errCommit = err
	// return
}

func (slf *pqInstanceTx) Rollback() (err error) {
	if slf == nil {
		return
	}

	if slf.isCommit || slf.isRollback || slf.tx == nil {
		if slf.errCommit != nil && slf.isCommit {
			log.Printf("db.tx.rollback cant commit: isCommit = %v\n", slf.isCommit)
		}

		if slf.isRollback {
			log.Printf("db.tx.rollback cant commit: isRollback = %v\n", slf.isRollback)
		}

		// if slf.tx == nil {
		// 	log.Printf("db.tx.rollback cant commit: tx = nil\n")
		// }

		return nil
	}

	iteration := 0
	const maxTry = 3

	for {
		iteration++
		err = slf.tx.Rollback()

		if err == nil {
			break
		}

		if err == sql.ErrTxDone {
			err = nil
			break
		}

		if iteration >= maxTry {
			break
		}

		time.Sleep(time.Microsecond * time.Duration(rand.Int63n(max-min)+min))
	}

	slf.isRollback = err == nil
	return

	// slf.isRollback = true
	// err = slf.tx.Rollback()
	// if err != nil {
	// 	log.Printf("db.tx.rollback error: %+v\n", err)
	// }
	// slf.errRollback = err

	// return
}

func (slf *pqInstanceTx) Host() string {
	return slf.instance.connRW.host
}
