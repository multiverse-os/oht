package network

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WS struct {
	Server *WebServer
	Client *websocket.Conn
	Engine *gin.Engine
}

func InitializeWS(listenURL *url.Url) *WS {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	ws := &WS{
		Server: InitializeWebServer(engine, listenURL.Host),
		Engine: engine,
	}
	engine.GET("/", func(c *gin.Context) {
		Serve(c.Writer, c.Request)
	})
	return webSocket
}

func (ws *WS) Listen(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (ws *WS) Stop() {
}

func (ws *WS) Connect(peerURL *url.URL) bool {
	d := websocket.Dialer{
		NetDial:          DialProxy(socks.SOCKS5, peerURL.Host),
		HandshakeTimeout: 15 * time.Second,
	}
	ws.Client, _, err = d.Dial(peerURL.String(), nil)
	return (err != nil)
}

func (ws *WebSocket) ReadMessages() {
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

func (ws *WS) WriteMessage(messageType int, payload types.Message) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteJSON(payload)
}

func (ws *WS) writeControl(messageType int) error {
	return ws.WriteControl(messageType, []byte{}, time.Now().Add(writeWait))
}

func (ws *WS) WriteMessages(p *Peer) {
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
