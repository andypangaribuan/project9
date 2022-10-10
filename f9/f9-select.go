/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

func RSelectFirst[T any](out []T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, nil
	}

	return &out[0], nil
}
