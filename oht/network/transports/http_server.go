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
		server.Manager.Serve(c.Writer, c.Request)
	})

	return websocket
}
