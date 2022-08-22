/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package crypto

type srCryptoMD5 struct{}

type srCryptoSHA256 struct{}

func Create() (*srCryptoMD5, *srCryptoSHA256) {
	return &srCryptoMD5{}, &srCryptoSHA256{}
}
