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
