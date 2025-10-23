/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"net/mail"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
	"github.com/golang-jwt/jwt/v5"
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

func (slf *srUtil) GetXID25() string {
	hex := fmt.Sprintf("%v", f9.TimeNow().UnixMicro())
	nine := slf.GetRandomAlphabetNumber(9)
	xid := hex + nine
	return xid
}

func (slf *srUtil) GetXID30() string {
	hex := fmt.Sprintf("%v", f9.TimeNow().UnixMicro())
	nine := slf.GetRandomAlphabetNumber(14)
	xid := hex + nine
	return xid
}

func (slf *srUtil) GetXID40() string {
	hex := fmt.Sprintf("%v", f9.TimeNow().UnixMicro())
	nine := slf.GetRandomAlphabetNumber(24)
	xid := hex + nine
	return xid
}

func (*srUtil) GetRandom(length int, value string) (string, error) {
	return gonanoid.Generate(value, length)
}

func (*srUtil) GetRandomNumber(min, max int) int {
	// rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (*srUtil) GetRandomAlphabet(length int) string {
	if length <= 0 {
		return ""
	}

	val, _ := gonanoid.Generate(idAlphabetLower+idAlphabetUpper, length)
	return val
}

func (*srUtil) GetRandomAlphabetLower(length int) string {
	if length <= 0 {
		return ""
	}

	val, _ := gonanoid.Generate(idAlphabetLower, length)
	return val
}

func (*srUtil) GetRandomAlphabetUpper(length int) string {
	if length <= 0 {
		return ""
	}

	val, _ := gonanoid.Generate(idAlphabetUpper, length)
	return val
}

func (*srUtil) GetRandomAlphabetNumber(length int) string {
	if length <= 0 {
		return ""
	}

	val, _ := gonanoid.Generate(idAlphabetLower+idAlphabetUpper+idNumeric, length)
	return val
}

func (*srUtil) GetRandomDuration(min int, max int, base time.Duration) time.Duration {
	x := rand.Int63n(int64(max)-int64(min)) + int64(min)
	return base * time.Duration(x)
}

func (slf *srUtil) BuildJwtToken(privateKey []byte, claims jwt.Claims) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (slf *srUtil) BuildJwtTokenWithPassword(privateKey []byte, password string, claims jwt.Claims) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(privateKey, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

// func (slf *srUtil) BuildJwtTokenWithPassword(privateKey []byte, password string, claims jwt.Claims) (string, error) {
// 	block, _ := pem.Decode(privateKey)
// 	if block == nil {
// 		return "", errors.New("invalid PEM: no block found")
// 	}

// 	var rsaKey *rsa.PrivateKey

// 	switch block.Type {
// 	case "ENCRYPTED PRIVATE KEY":
// 		if password == "" {
// 			return "", errors.New("password required for encrypted PKCS#8 key")
// 		}
// 		// Modern encrypted PKCS#8 (PBES2). No deprecated APIs used.
// 		keyAny, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, []byte(password))
// 		if err != nil {
// 			return "", fmt.Errorf("parse encrypted PKCS#8: %w", err)
// 		}
// 		k, ok := keyAny.(*rsa.PrivateKey)
// 		if !ok {
// 			return "", errors.New("encrypted key is not RSA")
// 		}
// 		rsaKey = k

// 	case "PRIVATE KEY":
// 		// Unencrypted PKCS#8
// 		keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 		if err != nil {
// 			return "", fmt.Errorf("parse PKCS#8: %w", err)
// 		}
// 		k, ok := keyAny.(*rsa.PrivateKey)
// 		if !ok {
// 			return "", errors.New("PKCS#8 key is not RSA")
// 		}
// 		rsaKey = k

// 	case "RSA PRIVATE KEY":
// 		// Unencrypted PKCS#1
// 		k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 		if err != nil {
// 			return "", fmt.Errorf("parse PKCS#1: %w", err)
// 		}
// 		rsaKey = k

// 	default:
// 		// Reject legacy PEM encryption (insecure by design) and other types.
// 		if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
// 			return "", errors.New(
// 				"legacy PEM encryption detected (Proc-Type: 4,ENCRYPTED) — unsupported; "+
// 					"convert to encrypted PKCS#8 (PBES2), e.g.: "+
// 					`openssl pkcs8 -topk8 -v2 aes-256-cbc -iter 600000 -in old.pem -out new.pem`,
// 			)
// 		}
// 		return "", fmt.Errorf("unsupported PEM type: %q", block.Type)
// 	}

// 	// Sign JWT (RS256)
// 	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	signed, err := tok.SignedString(rsaKey)
// 	if err != nil {
// 		return "", fmt.Errorf("sign jwt: %w", err)
// 	}
// 	return signed, nil
// }

func (slf *srUtil) CreateJwtToken(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	return slf.createJwtToken(key, subject, id, expiresAt, issuedAt, notBefore)
}

func (slf *srUtil) CreateJwtTokenWithPassword(subject, id string, expiresAt, issuedAt, notBefore time.Time, privateKey []byte, password string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(privateKey, password)
	if err != nil {
		return "", err
	}

	return slf.createJwtToken(key, subject, id, expiresAt, issuedAt, notBefore)
}

func (*srUtil) createJwtToken(key *rsa.PrivateKey, subject, id string, expiresAt, issuedAt, notBefore time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        id,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		NotBefore: jwt.NewNumericDate(notBefore),
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (*srUtil) GetJwtClaims(token string, publicKey []byte) (*jwt.RegisteredClaims, bool, error) {
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	// if err, ok := err.(*jwt.ValidationError); ok {
	// 	return claims, true, err
	// }

	if err != nil {
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

func (*srUtil) GetExecutionInfo(depth int) (execFunc string, execPath string) {
	pc, filename, line, _ := runtime.Caller(1 + depth)
	execFunc = runtime.FuncForPC(pc).Name()
	execPath = fmt.Sprintf("%v:%v", filename, line)
	return
}

func (*srUtil) IsNumberOnly(value string, exclude ...string) bool {
	v := value
	for _, e := range exclude {
		v = strings.ReplaceAll(v, e, "")
	}

	for i := 0; i < len(v); i++ {
		oneChar := v[i : i+1]
		_, err := strconv.Atoi(oneChar)
		if err != nil {
			return false
		}
	}

	return true
}

// country id:   https://laendercode.net/en/2-letter-list.html
// country code: https://countrycode.org/
func (*srUtil) ExtractPhoneNumber(phoneNumber *string) (countryId, countryCode, number string) {
	if phoneNumber == nil {
		return
	}

	v := *phoneNumber
	if v[:1] == "+" {
		v = v[1:]
	}

	if countryCode == "" && v[:1] == "0" {
		countryId = "ID"
		countryCode = "62"
		number = v[1:]
		return
	}

	if countryCode == "" {
		switch v[:2] {
		case "60": // malaysia
			countryId = "MY"
		case "61": // australia
			countryId = "AU"
		case "62": // indonesia
			countryId = "ID"
		case "63": // philippines
			countryId = "PH"
		case "65": // singapore
			countryId = "SG"
		case "66": // thailand
			countryId = "TH"
		}

		if countryId != "" {
			countryCode = v[:2]
			number = v[2:]
		}
	}

	return
}

// verifyDomain default false
func (*srUtil) IsEmailValid(email string, verifyDomain ...bool) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if len(verifyDomain) > 0 && verifyDomain[0] {
		parts := strings.Split(email, "@")
		mx, err := net.LookupMX(parts[1])
		if err != nil || len(mx) == 0 {
			return false
		}
	}

	return true
}
