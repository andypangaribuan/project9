/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package i9

import "github.com/andypangaribuan/project9/abs"

func GCS(bucketName string, credential []byte) (abs.I9GCS, error) {
	return newGcsInstance(bucketName, credential)
}
