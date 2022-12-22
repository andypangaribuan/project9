/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

/* spell-checker: disable */
package server

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

var cip *srClientIP

func init() {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	cidrs := make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}

	cip = &srClientIP{
		// Should use canonical format of the header key s
		// https://golang.org/pkg/net/http/#CanonicalHeaderKey

		// Header may return multiple IP addresses in the format: "client IP, proxy 1 IP, proxy 2 IP", so we take the the first one.
		cidrs:                       cidrs,
		xOriginalForwardedForHeader: http.CanonicalHeaderKey("X-Original-Forwarded-For"),
		xForwardedForHeader:         http.CanonicalHeaderKey("X-Forwarded-For"),
		xForwardedHeader:            http.CanonicalHeaderKey("X-Forwarded"),
		forwardedForHeader:          http.CanonicalHeaderKey("Forwarded-For"),
		forwardedHeader:             http.CanonicalHeaderKey("Forwarded"),

		// Standard headers used by Amazon EC2, Heroku, and others
		xClientIPHeader: http.CanonicalHeaderKey("X-Client-IP"),
		// Nginx proxy/FastCGI
		xRealIPHeader: http.CanonicalHeaderKey("X-Real-IP"),
		// Cloudflare.
		// @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
		// CF-Connecting-IP - applied to every request to the origin.
		cfConnectingIPHeader: http.CanonicalHeaderKey("CF-Connecting-IP"),
		// Fastly CDN and Firebase hosting header when forwared to a cloud function
		fastlyClientIPHeader: http.CanonicalHeaderKey("Fastly-Client-Ip"),
		// Akamai and Cloudflare
		trueClientIPHeader: http.CanonicalHeaderKey("True-Client-Ip"),
	}
}

func (slf *srClientIP) getClientIP(fuseCtx *srFuseContext) string {
	ctx := fuseCtx.fiberCtx.Context()

	xClientIP := ctx.Request.Header.Peek(slf.xClientIPHeader)
	if xClientIP != nil {
		return string(xClientIP)
	}

	xOriginalForwardedFor := ctx.Request.Header.Peek(slf.xOriginalForwardedForHeader)
	if xOriginalForwardedFor != nil {
		requestIP, err := slf.retrieveForwardedIP(string(xOriginalForwardedFor))
		if err == nil {
			return requestIP
		}
	}

	xForwardedFor := ctx.Request.Header.Peek(slf.xForwardedForHeader)
	if xForwardedFor != nil {
		requestIP, err := slf.retrieveForwardedIP(string(xForwardedFor))
		if err == nil {
			return requestIP
		}
	}

	if ip, err := slf.fromSpecialHeaders(fuseCtx); err == nil {
		return ip
	}

	if ip, err := slf.fromForwardedHeaders(fuseCtx); err == nil {
		return ip
	}

	var remoteIP string
	remoteAddr := ctx.RemoteAddr().String()

	if strings.ContainsRune(remoteAddr, ':') {
		remoteIP, _, _ = net.SplitHostPort(remoteAddr)
	} else {
		remoteIP = remoteAddr
	}
	return remoteIP
}

func (slf *srClientIP) fromSpecialHeaders(fuseCtx *srFuseContext) (string, error) {
	ctx := fuseCtx.fiberCtx.Context()

	ipHeaders := [...]string{slf.cfConnectingIPHeader, slf.fastlyClientIPHeader, slf.trueClientIPHeader, slf.xRealIPHeader}
	for _, iplHeader := range ipHeaders {
		if clientIP := ctx.Request.Header.Peek(iplHeader); clientIP != nil {
			return string(clientIP), nil
		}
	}
	return "", errors.New("can't get ip from special headers")
}

func (slf *srClientIP) fromForwardedHeaders(fuseCtx *srFuseContext) (string, error) {
	ctx := fuseCtx.fiberCtx.Context()

	forwardedHeaders := [...]string{slf.xForwardedHeader, slf.forwardedForHeader, slf.forwardedHeader}
	for _, forwardedHeader := range forwardedHeaders {
		if forwarded := ctx.Request.Header.Peek(forwardedHeader); forwarded != nil {
			if clientIP, err := slf.retrieveForwardedIP(string(forwarded)); err == nil {
				return clientIP, nil
			}
		}
	}
	return "", errors.New("can't get ip from forwarded headers")
}

// isLocalAddress works by checking if the address is under private CIDR blocks.
// List of private CIDR blocks can be seen on :
//
// https://en.wikipedia.org/wiki/Private_network
//
// https://en.wikipedia.org/wiki/Link-local_address
func (slf *srClientIP) isPrivateAddress(address string) (bool, error) {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}

	for i := range slf.cidrs {
		if slf.cidrs[i].Contains(ipAddress) {
			return true, nil
		}
	}

	return false, nil
}

func (slf *srClientIP) retrieveForwardedIP(forwardedHeader string) (string, error) {
	for _, address := range strings.Split(forwardedHeader, ",") {
		if len(address) > 0 {
			address = strings.TrimSpace(address)
			isPrivate, err := slf.isPrivateAddress(address)
			switch {
			case !isPrivate && err == nil:
				return address, nil
			case isPrivate && err == nil:
				return "", errors.New("forwarded ip is private")
			default:
				return "", err
			}
		}
	}
	return "", errors.New("empty or invalid forwarded header")
}
