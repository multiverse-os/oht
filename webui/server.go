package webui

import (
	"../network"
	"github.com/gin-gonic/gin"
	"strconv"
)

func InitializeServer(wsHost string, port int) {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.LoadHTMLFiles("webui/index.html")
	server.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"wsHost": wsHost,
		})
	})
	server.GET("/ws", func(c *gin.Context) {
		network.Manager.Serve(c.Writer, c.Request)
	})
	go server.Run("127.0.0.1:" + strconv.Itoa(port))
}
