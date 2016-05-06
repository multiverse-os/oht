package webui

import (
	"../../network"

	"github.com/gin-gonic/gin"
)

func InitializeServer(onionHost, webUIPort string) (server *network.WebServer) {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.LoadHTMLFiles("ui/webui/index.html")
	//r.Static("/public/css/", "./public/css")
	//r.Static("/public/js/", "./public/js/")
	//r.Static("/public/fonts/", "./public/fonts/")
	//r.Static("/public/img/", "./public/img/")

	//r.GET("/", IndexRouter)
	//r.GET("/about", AboutRoute)
	//r.GET("/contact", ContactRoute)
	//r.GET("/signin", SigninRoute)
	//r.GET("/signup", SignupRoute)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"wsHost": onionHost,
		})
	})
	return network.InitializeWebServer(r, (onionHost + ":" + webUIPort))
}
