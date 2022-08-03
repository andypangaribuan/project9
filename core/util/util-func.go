/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"time"

	"github.com/andypangaribuan/project9/p9"
	"github.com/golang-jwt/jwt"
	"github.com/matoous/go-nanoid/v2"
)

func (*srUtil) GetNanoID(length ...int) (string, error) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	size := p9.Conf.NanoIdLength
	if len(length) > 0 {
		size = length[0]
	}

	return gonanoid.Generate(alphabet, size)
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
