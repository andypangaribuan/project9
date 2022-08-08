/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
	"github.com/golang-jwt/jwt"
	"github.com/matoous/go-nanoid/v2"
)

const idAlphabetLower = "abcdefghijklmnopqrstuvwxyz"
const idAlphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const idNumeric = "0123456789"

func (slf *srUtil) GetNanoID(length ...int) string {
	size := p9.Conf.NanoIdLength
	if len(length) > 0 {
		size = length[0]
	}
	if size <= 0 {
		return ""
	}

	id, _ := slf.GetRandom(size, idNumeric+idAlphabetLower+idAlphabetUpper)
	return id
}

func (slf *srUtil) GetID25() string {
	unixMicro := fmt.Sprintf("%v", f9.TimeNow().UnixMicro())
	u3 := unixMicro[len(unixMicro)-3:]
	ul := unixMicro[:len(unixMicro)-3]
	id := fmt.Sprintf("%v%v", u3, ul)

	nn, _ := slf.GetRandom(9, idAlphabetLower+idAlphabetUpper)
	n1 := nn[:6]
	n2 := nn[6:]

	return fmt.Sprintf("%v%v%v", n1, id, n2)
}

func (*srUtil) GetRandom(length int, value string) (string, error) {
	return gonanoid.Generate(value, length)
}

func (*srUtil) CreateJwtToken(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Id:        id,
		IssuedAt:  issuedAt.Unix(),
		NotBefore: notBefore.Unix(),
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (*srUtil) GetJwtClaims(token string, publicKey []byte) (*jwt.StandardClaims, bool, error) {
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims := &jwt.StandardClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	if err, ok := err.(*jwt.ValidationError); ok {
		return claims, true, err
	}

	return claims, false, err
}

func (*srUtil) Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (*srUtil) Base64Decode(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}
