package p2p

import (
	"net"
	"net/url"
	"time"

	"../../network"
	"../../types"
)

type TCPClient struct {
	Client     *net.Conn
	onionhost  string
	remotehost string
	socksPort  string
}

func (tcpClient *TCPClient) Connect(remotehost, socksPort string) bool {
	//u := url.URL{Scheme: "tcp", Host: remotehost, Path: "/"}
	//NetDial:          network.DialProxy(socks.SOCKS5, ("127.0.0.1:" + socksPort)),
	return false
}
