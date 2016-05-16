package transports

import (
	"net"
	"net/url"
	"time"

	"../../network"
)

type HTTPClient struct {
	//Client     *websocket.Conn
	onionhost     string
	remotehost    string
	socksPort     string
	routingConfig string // routing defined by a json config
}

func (hClient *HTTPClient) Connect(remotehost, socksPort string) bool {
	err := nil
	if err != nil {
		return false
	} else {
		return true
	}
}

func (hClient *HTTPClient) readMessages() {
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
