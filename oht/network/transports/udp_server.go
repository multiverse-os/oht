package transports

import (
	"../../network"
)

// Note: Will be tunneled over TCP

type UDPServer struct {
	Server    string
	Engine    string
	Onionhost string
}

func InitializeUDPServer(onionhost, udpPort string) *UDPServer {
	udp := &UDPServer{
		Server:    "udp:",
		Onionhost: onionhost,
	}

	return udp
}

func Start() {
}

func Stop() {
}
