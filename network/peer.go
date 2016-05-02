package network

import (
	"log"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Peer struct {
	// Each peer needs its own crypto key to encrypt and decrypt shit
	Id         string
	Version    int
	Reputation int32
	Ignored    bool
	OnionHost  string

	Websocket *websocket.Conn
	Send      chan Message
	Data      interface{}
}

type Message struct {
	Type       string
	Id         string
	Timestamp  int64  `json:",omitempty"`
	OriginHost string `json:",omitempty"`
	Username   string `json:",omitempty"`
	Body       string `json:",omitempty"`
}

const (
	writeWait  = 15 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectToPeer(peerHost string, socksPort string) {
	sock := &Socks4a{Network: "tcp", Address: "127.0.0.1:" + socksPort}
	u := url.URL{Scheme: "ws", Host: peerHost, Path: "/ws"}
	d := websocket.Dialer{
		NetDial:          func(network, addr string) (net.Conn, error) { return sock.Dial(peerHost) },
		HandshakeTimeout: 15 * time.Second,
	}
	ws, _, err := d.Dial(u.String(), nil)
	if err != nil {
		log.Println("Failed to connect to peer: ", err)
	} else {
		p := &Peer{
			Send:       make(chan Message, 256),
			Websocket:  ws,
			OnionHost:  peerHost,
			Version:    1,
			Reputation: 0,
			Ignored:    false,
		}
		Manager.Register <- p
		go p.writeMessages()
		p.readMessages()
	}
}

func (p *Peer) readMessages() {
	defer func() {
		Manager.Unregister <- p
	}()
	p.Websocket.SetReadLimit(maxMessageSize)
	p.Websocket.SetReadDeadline(time.Now().Add(pongWait))
	p.Websocket.SetPongHandler(func(string) error {
		p.Websocket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var message Message
		err := p.Websocket.ReadJSON(&message)
		if err != nil {
			break
		}
		Manager.Receive <- message
	}
}

func (p *Peer) writeMessage(messageType int, payload Message) error {
	p.Websocket.SetWriteDeadline(time.Now().Add(writeWait))
	return p.Websocket.WriteJSON(payload)
}

func (p *Peer) writeControl(messageType int) error {
	return p.Websocket.WriteControl(messageType, []byte{}, time.Now().Add(writeWait))
}

func (p *Peer) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		Manager.Unregister <- p
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-p.Send:
			if !ok {
				p.writeControl(websocket.CloseMessage)
				return
			}
			_ = p.writeMessage(websocket.TextMessage, message)
		case <-ticker.C:
			_ = p.writeControl(websocket.PingMessage)
		}
	}
}
