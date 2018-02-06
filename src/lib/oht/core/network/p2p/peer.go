package network

import (
	"github.com/gorilla/websocket"
)

type OnionServiceConfig struct {
	DirectoryName   string
	OnionHost       string
	OnionPrivateKey string
	ListenPort      string
}

type Peer struct {
	Id        string
	Config    *OnionServiceConfig
	Connected int8
	WebSocket *websocket.Conn
	Manager   *Manager
	Send      chan Message
	Data      interface{}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *Peer) readMessages() {
}

func (p *Peer) writeMessage(message Message) error {
	return nil
}

func (p *Peer) writeControl(messageType int) error {
	return nil
}

func (p *Peer) writeMessages() {
}
