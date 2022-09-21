/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package crypto

type srCryptoAES struct{}

type srCryptoMD5 struct{}

type srCryptoSHA256 struct{}

func Create() (*srCryptoAES, *srCryptoMD5, *srCryptoSHA256) {
	return &srCryptoAES{}, &srCryptoMD5{}, &srCryptoSHA256{}
}
