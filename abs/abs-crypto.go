/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

// https://www.devglan.com/online-tools/aes-encryption-decryption
type CryptoAES interface {
	ECBEncrypt(key, value string) (string, error)
	ECBRawEncrypt(key, value string) ([]byte, error)
	ECBDecrypt(key, value string) (string, error)
	ECBRawDecrypt(key string, value []byte) ([]byte, error)
}

type CryptoMD5 interface {
	Generate(data string) string
}

type CryptoSHA256 interface {
	Generate(data string) string
}

type CryptoSHA512 interface {
	Generate(data string) string
}
