/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package check

import "strings"

func (*srStr) IsEmptyPtr(val *string) (string, bool) {
	if val == nil {
		return "", true
	}

	res := strings.TrimSpace(*val)
	if res == "" {
		return "", true
	}

	return res, false
}
