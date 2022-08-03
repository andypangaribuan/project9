/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package crypto

type srCrypto struct{}

func Create() *srCrypto {
	return &srCrypto{}
}
