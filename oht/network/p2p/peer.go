package p2p

import (
	"net"
	"net/url"
	"time"

	"../../network"
	"../../types"

	"github.com/gorilla/websocket"
)

type Peer struct {
	// Each peer needs its own crypto key to encrypt and decrypt shit
	Id         string
	Version    int
	Reputation int32
	Ignored    bool
	OnionHost  string
	WebSocket  *websocket.Conn
	Manager    *Manager
	Send       chan types.Message
	Data       interface{}
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

func (manager *Manager) ConnectToPeer(peerHost, socksPort string) bool {
	sock := &network.Socks4a{Network: "tcp", Address: ("127.0.0.1:" + socksPort)}
	u := url.URL{Scheme: "ws", Host: peerHost, Path: "/ws"}
	d := websocket.Dialer{
		NetDial:          func(network, addr string) (net.Conn, error) { return sock.Dial(peerHost) },
		HandshakeTimeout: 15 * time.Second,
	}
	ws, _, err := d.Dial(u.String(), nil)
	if err != nil {
		return false
	} else {
		p := &Peer{
			Send:       make(chan types.Message, 256),
			WebSocket:  ws,
			Manager:    manager,
			OnionHost:  peerHost,
			Version:    1,
			Reputation: 0,
			Ignored:    false,
		}
		manager.Register <- p
		go p.writeMessages()
		p.readMessages()
		return true
	}
}

func (p *Peer) readMessages() {
	defer func() {
		p.Manager.Unregister <- p
	}()
	p.WebSocket.SetReadLimit(maxMessageSize)
	p.WebSocket.SetReadDeadline(time.Now().Add(pongWait))
	p.WebSocket.SetPongHandler(func(string) error {
		p.WebSocket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var message types.Message
		err := p.WebSocket.ReadJSON(&message)
		if err != nil {
			break
		}
		p.Manager.Receive <- message
	}
}

func (p *Peer) writeMessage(messageType int, payload types.Message) error {
	p.WebSocket.SetWriteDeadline(time.Now().Add(writeWait))
	return p.WebSocket.WriteJSON(payload)
}

func (p *Peer) writeControl(messageType int) error {
	return p.WebSocket.WriteControl(messageType, []byte{}, time.Now().Add(writeWait))
}

func (p *Peer) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		p.Manager.Unregister <- p
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
