/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type srCrypto struct {
	AES    *srCryptoAES
	MD5    *srCryptoMD5
	SHA256 *srCryptoSHA256
	SHA512 *srCryptoSHA512
}

type srCryptoAES struct {
	absCryptoAES
}

type absCryptoAES interface {
	abs.CryptoAES
}

type srCryptoMD5 struct {
	absCryptoMD5
}

type absCryptoMD5 interface {
	abs.CryptoMD5
}

type srCryptoSHA256 struct {
	absCryptoSHA256
}

type absCryptoSHA256 interface {
	abs.CryptoSHA256
}

type srCryptoSHA512 struct {
	absCryptoSHA512
}

type absCryptoSHA512 interface {
	abs.CryptoSHA512
}