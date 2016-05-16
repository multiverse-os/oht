package network

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Client     *websocket.Conn
	onionhost  string
	remotehost string
	socksPort  string
}

func (wsClient *WebSocketClient) Connect(remotehost, socksPort string) bool {
	u := url.URL{Scheme: "ws", Host: remotehost, Path: "/"}
	d := websocket.Dialer{
		NetDial:          DialProxy(socks.SOCKS5, ("127.0.0.1:" + socksPort)),
		HandshakeTimeout: 15 * time.Second,
	}
	ws, _, err := d.Dial(u.String(), nil)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (wsClient *WebSocketClient) readMessages() {
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

func (ws *WebSocketClient) writeMessage(messageType int, payload types.Message) error {
	ws.WebSocket.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WebSocket.WriteJSON(payload)
}

func (ws *WebSocketClient) writeControl(messageType int) error {
	return ws.WebSocket.WriteControl(messageType, []byte{}, time.Now().Add(writeWait))
}

func (ws *WebSocketClient) writeMessages(p *Peer) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ws.Manager.Unregister <- p
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-p.Send:
			if !ok {
				p.writeControl(websocket.CloseMessage)
				return
			}
			_ = ws.writeMessage(websocket.TextMessage, message)
		case <-ticker.C:
			_ = ws.writeControl(websocket.PingMessage)
		}
	}
}
