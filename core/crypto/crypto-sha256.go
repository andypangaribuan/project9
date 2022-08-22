/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func (*srCryptoSHA256) Generate(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
