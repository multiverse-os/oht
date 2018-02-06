package network

import ()

type TCP struct {
	ListenURL *url.URL
	Engine    string // will be net.Conn
}

func InitializeTCP(listenURL *url.URL) *TCP {
	return &TCP{
		ListenURL: listenURL,
	}
}

func Listen() {}
func Stop() {}
//func 
