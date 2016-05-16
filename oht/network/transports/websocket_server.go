package transports

import (
	"../../network"
	"github.com/gin-gonic/gin"
)

type WebsocketServer struct {
	Server    *network.WebServer
	Engine    *gin.Engine
	Onionhost string
}

type Peer struct {
	Id        string
	Connected int8
	WebSocket *websocket.Conn
	Manager   *network.Manager
	Send      chan types.Message
	Data      interface{}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func InitializeWebsocket(onionhost, websocketPort string) *WebsocketServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	websocket := &WebsocketServer{
		Server:    network.InitializeWebServer(engine, ("127.0.0.1:" + websocketPort)),
		Engine:    engine,
		Onionhost: onionhost,
	}

	engine.GET("/", func(c *gin.Context) {
		Serve(c.Writer, c.Request)
	})

	return websocket
}

func (manager *Manager) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	p := &Peer{Send: make(chan types.Message, 256), WebSocket: ws}
	manager.Register <- p
	go p.writeMessages()
	p.readMessages()
}
