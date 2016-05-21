package network

import (
	"log"
	"net/url"

	"../network/p2p"
	"../network/transports"
)

type ConnectionPool struct {
	Connections []*Connection
}

type Connection struct {
	Transport interface{}
	Protocol  string
	ListenURL *url.URL
}

type Transport interface {
	Listen(listenURL *url.URL)
	Stop()
	Connect(peerURL *url.URL)
	Read(peerURL *url.URL, message string)
	Write(peerURL *url.URL, message string)
}

func InitializeConnectionPool() *ConnectionPool {
	return &Connections{}
}

func (cp *ConnectionPool) InitializeConnection(listenURL) *ConnectionPool {
	connection := &Connection{ListenURL: listenURL}
	// Detect subprotocols
	if connection.Scheme == "oht" {
		connection.Protocol = "oht"
		connection.ListenURL.Scheme = "tcp"
	} else if connection.Scheme == "ricochet" {
		connection.Protocol = "ricochet"
		connection.ListenURL.Scheme = "tcp"
	} else if len(connection.Path) > 1 {
		connection.Protocol = connection.Path[1:]
	}
	// Initialize Transport
	if connection.Scheme == "tcp" {
		connection.Transport = InitializeTCP(connection.listenURL)
	} else if connection.Scheme == "http" {
		connection.Transport = InitializeHTTP(connection.listenURL)
	} else if connection.Scheme == "ws" {
		connection.Transport = InitializeWS(connection.listenURL)
	}
	return cp.Append(collection)
}

func (cp *ConnectionPool) Interface(connection Transport, action string, remoteURL *url.URL) {
	if action == "listen" {
		connection.Listen(connection.listenURL)
	} else if action == "stop" {
		connection.Stop()
	} else if action == "connect" {
		connection.Connect(remoteURL)
		p := &Peer{Connected: true, Send: make(chan Message, 256), Connection: connection, LastSeen: time.Time()}
		ws.Manager.Register <- p
		go ws.writeMessages()
		ws.readMessages()
	} else if action == "read" {
		connection.Read("peerAddress")
	} else if action == "write" {
		connection.Write("peerAddress")
	}
}
