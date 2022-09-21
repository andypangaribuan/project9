/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import "strings"

 func PtrStringTrimSpace(pointers ...*string) {
	for i := 0; i < len(pointers); i++ {
		if pointers[i] != nil {
			val := strings.TrimSpace(*pointers[i])
			*pointers[i] = val
		}
	}
}
