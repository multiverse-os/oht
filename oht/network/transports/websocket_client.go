package p2p

import (
	"net"
	"net/url"
	"time"

	"../../network"
	"../../types"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	Client    *websocket.Conn
	OnionHost string
}

func (wsClient *WebsocketClient) Connect(remotehost, socksPort string) bool {
	u := url.URL{Scheme: "ws", Host: remotehost, Path: "/"}
	d := websocket.Dialer{
		NetDial:          network.DialProxy(socks.SOCKS5, ("127.0.0.1:" + socksPort)),
		HandshakeTimeout: 15 * time.Second,
	}
	ws, _, err := d.Dial(u.String(), nil)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (wsClient *WebsocketClient) readMessages() {
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
