/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

type CheckStr interface {
	IsEmpty(val string) bool
	IsEmptyPtr(val *string) (string, bool)
}
