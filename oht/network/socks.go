// Copyright 2012, Hailiang Wang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

const (
	SOCKS4 = iota
	SOCKS4A
	SOCKS5
)

func DialProxy(socksType int, proxy string) func(string, string) (net.Conn, error) {
	if socksType == SOCKS5 {
		return func(_, targetAddr string) (conn net.Conn, err error) {
			return dialSocks5(proxy, targetAddr)
		}
	}
	return func(_, targetAddr string) (conn net.Conn, err error) {
		return dialSocks4(socksType, proxy, targetAddr)
	}
}

func dialSocks5(proxy, targetAddr string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", proxy)
	if err != nil {
		return
	}
	// version identifier/method selection request
	req := []byte{
		5, // version number
		1, // number of methods
		0, // method 0: no authentication (only anonymous access supported for now)
	}
	resp, err := sendReceive(conn, req)
	if err != nil {
		return
	} else if len(resp) != 2 {
		err = errors.New("Proxy: Server did not respond correctly.")
		return
	} else if resp[0] != 5 {
		err = errors.New("Proxy: Server does not support Socks 5.")
		return
	} else if resp[1] != 0 { // no auth
		err = errors.New("Proxy: Negotiation failed.")
		return
	}
	// detail request
	host, port, err := splitHostPort(targetAddr)
	req = []byte{
		5,               // version number
		1,               // connect command
		0,               // reserved, must be zero
		3,               // address type, 3 means domain name
		byte(len(host)), // address length
	}
	req = append(req, []byte(host)...)
	req = append(req, []byte{
		byte(port >> 8), // higher byte of destination port
		byte(port),      // lower byte of destination port (big endian)
	}...)
	resp, err = sendReceive(conn, req)
	if err != nil {
		return
	} else if len(resp) != 10 {
		err = errors.New("Proxy: Server did not respond correctly.")
	} else if resp[1] != 0 {
		err = errors.New("Proxy: Can't complete SOCKS5 connection.")
	}

	return
}

func dialSocks4(socksType int, proxy, targetAddr string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", proxy)
	if err != nil {
		return
	}

	host, port, err := splitHostPort(targetAddr)
	if err != nil {
		return
	}
	ip := net.IPv4(0, 0, 0, 1).To4()
	req := []byte{
		4,                          // version number
		1,                          // command CONNECT
		byte(port >> 8),            // higher byte of destination port
		byte(port),                 // lower byte of destination port (big endian)
		ip[0], ip[1], ip[2], ip[3], // special invalid IP address to indicate the host name is provided
		0, // user id is empty, anonymous proxy only
	}
	if socksType == SOCKS4A {
		req = append(req, []byte(host+"\x00")...)
	}

	resp, err := sendReceive(conn, req)
	if err != nil {
		return
	} else if len(resp) != 8 {
		err = errors.New("Proxy: Server did not respond correctly.")
		return
	}
	switch resp[1] {
	case 90:
		// request granted
	case 91:
		err = errors.New("Proxy: Socks connection request rejected or failed.")
	case 92:
		err = errors.New("Proxy: Socks connection request rejected becasue SOCKS server cannot connect to identd on the client.")
	case 93:
		err = errors.New("Proxy: Socks connection request rejected because the client program and identd report different user-ids.")
	default:
		err = errors.New("Proxy: Socks connection request failed, unknown error.")
	}
	return
}

func sendReceive(conn net.Conn, req []byte) (resp []byte, err error) {
	_, err = conn.Write(req)
	if err != nil {
		return
	}
	resp, err = readAll(conn)
	return
}

func readAll(conn net.Conn) (resp []byte, err error) {
	resp = make([]byte, 1024)
	n, err := conn.Read(resp)
	resp = resp[:n]
	return
}
