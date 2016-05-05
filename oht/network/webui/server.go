package webui

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	Port   string
}

func InitializeServer(onionHost, webUIPort string) (server *Server) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLFiles("ui/webui/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"wsHost": onionHost,
		})
	})
	return &Server{
		Router: router,
		Port:   webUIPort,
	}
}

func (server *Server) Start() bool {
	go server.Router.Run(":" + server.Port)
	return true
}

func (server *Server) Stop() bool {
	// Apparently gin does not have the ability to stop the server
	// another solution will have to be found
	return false
}
