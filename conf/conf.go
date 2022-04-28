/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conf

var (
	HashIdSalt   string
	HashIdLength int
)

func init() {
	HashIdSalt = "Project9"
	HashIdLength = 20
}
