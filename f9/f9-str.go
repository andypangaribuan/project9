/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import "strings"

func PtrStringTrimSpace(pointers ...*string) *string {
	if len(pointers) > 0 && pointers[0] != nil {
		ptrValue := pointers[0]
		value := strings.TrimSpace(*ptrValue)
		return &value
	}

	return nil
}
