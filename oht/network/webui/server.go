package webui

import (
	"../../network"

	"github.com/gin-gonic/gin"
)

var (
	host string
)

func InitializeServer(onionHost, webUIPort string) (server *network.WebServer) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	host = onionHost

	engine.LoadHTMLGlob("ui/webui/templates/*")
	engine.Static("/public/css/", "ui/webui/public/css")
	engine.Static("/public/js/", "ui/webui/public/js/")
	engine.Static("/public/fonts/", "ui/webui/public/fonts/")
	engine.Static("/public/img/", "ui/webui/public/img/")

	engine.GET("/", getIndex)
	engine.GET("/about", getAbout)
	engine.GET("/contact", getContact)

	return network.InitializeWebServer(engine, ("127.0.0.1:" + webUIPort))
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"wsHost": host,
	})
}

func getAbout(c *gin.Context) {
	c.HTML(200, "about.html", gin.H{
		"wsHost": host,
	})
}

func getContact(c *gin.Context) {
	c.HTML(200, "contact.html", gin.H{
		"wsHost": host,
	})
}
