/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package p9

import "github.com/andypangaribuan/project9/abs"

type srCrypto struct {
	MD5    *srCryptoMD5
	SHA256 *srCryptoSHA256
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
