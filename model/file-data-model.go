/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package model

type FileData struct {
	Name   string
	Data   []byte
	Width  *int
	Height *int
}
