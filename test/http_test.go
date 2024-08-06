/*
 * Copyright (c) 2024.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package test

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	// "github.com/andypangaribuan/project9"
	"github.com/andypangaribuan/project9/p9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttp(t *testing.T) {
	// project9.Initialize()

	// url := "https://connect.partner.treasury.id/private/client-ip"
	url := "https://api.treasury.id/api/v1/private/svc-version"
	data, code, err := p9.Http.Get(url, nil, nil, true, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	fmt.Printf("data: %v\n", string(data))
}

func TestHttpSSL(t *testing.T) {
	// trsCert := getCertificate("api.treasury.id")
	// println(trsCert)

	// _, err := http.Get("https://golang.org/")
	// require.Nil(t, err)

	// certFilePath := getDirPath() + "/cert.pem"
	certFilePath := getDirPath() + "/treasury-id-chain.pem"
	caCert, err := os.ReadFile(certFilePath)
	require.Nil(t, err)

	caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		require.Nil(t, errors.New("invalid certs"))
	}

	// for block, rest := pem.Decode(caCert); block != nil; block, rest = pem.Decode(rest) {
	//     switch block.Type {
	//     case "CERTIFICATE":
	//         cert, err := x509.ParseCertificate(block.Bytes)
	//         if err != nil {
	//             panic(err)
	//         }

	//         // Handle certificate
	//         fmt.Printf("%T %#v\n", cert, cert)

	//     case "PRIVATE KEY":
	//         key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	//         if err != nil {
	//             panic(err)
	//         }

	//         // Handle private key
	//         fmt.Printf("%T %#v\n", key, key)

	//     default:
	//         panic("unknown block type")
	//     }
	// }

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// res, err := client.Get("https://ipecho.net/plain")
	res, err := client.Get("https://api.treasury.id/api/v1/private/svc-version")
	// res, err := client.Get("https://connect.partner.treasury.id/private/client-ip")
	if err != nil {
		x, ok := err.(net.Error)
		if ok {
			log.Println(x)
		}
	}
	require.Nil(t, err)
	require.NotNil(t, res)
}

// func getCertificate(domain string) string {
//     conn, err := tls.Dial("tcp", fmt.Sprintf("%v:443", domain), nil)
//     if err != nil {
//         panic(err.Error())
//     }
//     defer conn.Close()
//     certificate := conn.ConnectionState().PeerCertificates[0]

// 	// x509.MarshalPKCS1PublicKey()
//     x509AsBytes, err := x509.MarshalPKIXPublicKey(certificate.PublicKey)
//     // x509AsBytes, err := x509.MarshalPKIXPublicKey(certificate.Raw)
//     if err != nil {
//         panic(err.Error())
//     }

//     pem := pem.EncodeToMemory(&pem.Block{
//         Type:  "CERTIFICATE",
//         Bytes: x509AsBytes,
//     })

//     return string(pem)
// }
