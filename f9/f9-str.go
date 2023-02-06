/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import "strings"

func PtrStringTrimSpace(pointers ...*string) *string {
	// for i := 0; i < len(pointers); i++ {
	// 	if pointers[i] != nil {
	// 		val := strings.TrimSpace(*pointers[i])
	// 		*pointers[i] = val
	// 	}
	// }

	if len(pointers) > 0 && pointers[0] != nil {
		ptrValue := pointers[0]
		value := strings.TrimSpace(*ptrValue)
		return &value
	}

	return nil
}
