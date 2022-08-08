/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package err

type srErr struct{}

func Create() *srErr {
	return &srErr{}
}
