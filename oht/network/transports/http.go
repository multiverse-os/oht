package network

import (
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	Server    *network.WebServer
	Engine    *gin.Engine
	Onionhost string
}

func InitializeHTTP(onionhost, httpPort string) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	http := &HTTPServer{
		Server:    network.InitializeWebServer(engine, ("127.0.0.1:" + httpPort)),
		Engine:    engine,
		Onionhost: onionhost,
	}

	engine.GET("/", func(c *gin.Context) {
		// Provide basic protocol features over REST
	})

	return http
}
