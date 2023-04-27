/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func NewGrpcClient[T any](address string, usingClientLB bool, newClientFunc func(*grpc.ClientConn) T) (*GrpcClient[T], error) {
	client := new(GrpcClient[T])
	client.address = address
	client.usingClientLB = usingClientLB
	client.newClientFunc = newClientFunc

	client.usingTLS = false
	arr := strings.Split(address, ":")
	if len(arr) > 1 {
		if arr[len(arr)-1] == "443" {
			client.usingTLS = true
		}
	}

	err := client.buildConnection()
	if err != nil {
		return nil, err
	}

	return client, err
}

func (slf *GrpcClient[T]) buildConnection() error {
	var (
		address           = slf.address
		grpcServiceConfig = ""
		conn              *grpc.ClientConn
		err               error
	)

	if slf.usingClientLB {
		resolver.SetDefaultScheme("dns")
		grpcServiceConfig = `{"loadBalancingPolicy":"round_robin"}`
		address = fmt.Sprintf("dns:///%v", address)
	}

	if slf.usingTLS {
		creds := credentials.NewTLS(&tls.Config{})
		if slf.usingClientLB {
			conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithDefaultServiceConfig(grpcServiceConfig))
		} else {
			conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
		}
	} else {
		if slf.usingClientLB {
			conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(grpcServiceConfig))
		} else {
			conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
	}

	if err != nil {
		return err
	}

	// make sure we connected to the service
	timeout := 5 * time.Second
	tcpConn, err := net.DialTimeout("tcp", slf.address, timeout)
	if err != nil {
		return err
	}
	_ = tcpConn.Close()

	slf.conn = conn
	slf.client = slf.newClientFunc(slf.conn)

	return err
}

func (slf *GrpcClient[T]) Client() T {
	return slf.client
}
