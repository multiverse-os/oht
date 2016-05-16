package transports

import ()

type TCPServer struct {
	Server    string
	Engine    string
	Onionhost string
}

func InitializeTCPServer(onionhost, tcpPort string) *TCPServer {
	tcp := &TCPServer{
		Server:    "tcp:",
		Onionhost: onionhost,
	}

	return tcp
}

func Start() {
}

func Stop() {
}
