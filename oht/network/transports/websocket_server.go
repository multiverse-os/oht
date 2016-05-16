package network

import (
	"time"

	"../p2p"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	Server    *WebServer
	Engine    *gin.Engine
	Onionhost string
}

type Peer struct {
	Id        string
	Connected bool
	WebSocket *websocket.Conn
	Manager   *Manager
	Send      chan Message
	Data      interface{}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func InitializeWebSocket(onionhost, webSocketPort string) *WebSocketServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	webSocket := &WebsocketServer{
		Server:    InitializeWebServer(engine, ("127.0.0.1:" + webSocketPort)),
		Engine:    engine,
		Onionhost: onionhost,
	}

	engine.GET("/", func(c *gin.Context) {
		Serve(c.Writer, c.Request)
	})
	return webSocket
}

func (ws *WebSocketServer) readMessages() {
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

func (ws *WebSocketServer) writeMessage(messageType int, payload types.Message) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteJSON(payload)
}

func (ws *WebSocketServer) writeControl(messageType int) error {
	return ws.WriteControl(messageType, []byte{}, time.Now().Add(writeWait))
}

func (ws *WebSocketServer) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	p := &Peer{Connected: true, Send: make(chan types.Message, 256), Connection: connection, LastSeen: time.Time()}
	ws.Manager.Register <- p
	go ws.writeMessages()
	ws.readMessages()
}
