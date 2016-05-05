package webui

import (
	"strconv"

	"./../../oht/network"

	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router  *gin.Engine
	Handler *mannders.GracefulServer
}

func InitializeServer() (server *Server) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLFiles("ui/webui/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"wsHost": tor.WebUIOnionHost,
		})
	})
	router.GET("/ws", func(c *gin.Context) {
		network.Manager.Serve(c.Writer, c.Request)
	})
	return &Server{
		Router: router,
	}
	server.Handler = manners
}

func (server *Server) Start() bool {
	err := server.Handler.ListenAndServe((":" + tor.WebUIPort), server.Router)
	return (err == nil)
}

func (server *Server) Stop() bool {
	return server.Handler.Close()
}
