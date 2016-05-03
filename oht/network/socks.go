package network

import (
	"errors"
	"io"
	"net"
	"strconv"
)

const (
	socks4aVersion       = 4
	socks4aConnect       = 1
	socks4aGranted       = 90
	socks4aRejected      = 91
	socks4aMissingIdentd = 92
	socks4aFailedIdentd  = 93
)

type Socks4a struct {
	Network string
	Address string
}

func (s *Socks4a) Dial(destination string) (net.Conn, error) {
	destStr, portStr, err := net.SplitHostPort(destination)
	if err != nil {
		return nil, err
	}
	dest := []byte(destStr)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("Proxy: Failed to parse port number: " + portStr)
	}
	if port < 1 || port > 0xffff {
		return nil, errors.New("Proxy: Port number out of range: " + portStr)
	}

	buf := make([]byte, 11+len(dest))
	buf[0] = socks4aVersion
	buf[1] = socks4aConnect
	buf[2] = byte(port >> 8)
	buf[3] = byte(port)
	buf[4] = 0
	buf[5] = 0
	buf[6] = 0
	buf[7] = 1
	buf[8] = 65
	buf[9] = 0
	for i, c := range dest {
		buf[10+i] = c
	}
	buf[10+len(dest)] = 0

	conn, err := net.Dial(s.Network, s.Address)
	if err != nil {
		return nil, err
	}

	closeConn := &conn
	defer func() {
		if closeConn != nil {
			(*closeConn).Close()
		}
	}()

	_, err = conn.Write(buf)
	if err != nil {
		return nil, errors.New("Proxy: Failed to write CONNECT message to " + s.Address + ": " + err.Error())
	}

	_, err = io.ReadFull(conn, buf[:8])
	if err != nil {
		return nil, errors.New("Proxy: Failed to connect to " + s.Address + ": " + err.Error())
	}

	if buf[0] != 0 {
		return nil, errors.New("Proxy: " + s.Address + " has unexpected version " + strconv.Itoa(int(buf[0])))
	}
	if buf[1] != socks4aGranted {
		return nil, errors.New("Proxy: " + s.Address + " failed acceptance")
	}
	closeConn = nil
	return conn, nil
}
