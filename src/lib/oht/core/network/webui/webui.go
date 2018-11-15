package webui

import (
	"html/template"

	"github.com/multiverse-os/libs/oht/core/network"

	"github.com/gin-gonic/gin"
)

type WebUI struct {
	Server       *network.WebServer
	Engine       *gin.Engine
	Onionhost    string
	BaseTemplate string
	Templates    map[string]*template.Template
}

func InitializeWebUI(onionhost, webUIPort string) *WebUI {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	webUI := &WebUI{
		Server:       network.InitializeWebServer(engine, ("127.0.0.1:" + webUIPort)),
		Engine:       engine,
		Onionhost:    onionhost,
		BaseTemplate: "ui/webui/templates/layouts/application.html",
		Templates:    make(map[string]*template.Template),
	}

	webUI.Templates["index"] = template.Must(template.ParseFiles(webUI.BaseTemplate, "ui/webui/templates/index.html"))
	webUI.Templates["about"] = template.Must(template.ParseFiles(webUI.BaseTemplate, "ui/webui/templates/about.html"))
	webUI.Templates["contact"] = template.Must(template.ParseFiles(webUI.BaseTemplate, "ui/webui/templates/contact.html"))

	engine.Static("/public/css/", "ui/webui/public/css")
	engine.Static("/public/js/", "ui/webui/public/js/")
	engine.Static("/public/fonts/", "ui/webui/public/fonts/")
	engine.Static("/public/img/", "ui/webui/public/img/")

	engine.GET("/", webUI.getIndex)
	engine.GET("/about", webUI.getAbout)
	engine.GET("/contact", webUI.getContact)

	return webUI
}

func (webUI *WebUI) getIndex(c *gin.Context) {
	webUI.Engine.SetHTMLTemplate(webUI.Templates["index"])
	c.HTML(200, "application.html", gin.H{
		"wsHost": webUI.Onionhost,
	})
}

func (webUI *WebUI) getAbout(c *gin.Context) {
	webUI.Engine.SetHTMLTemplate(webUI.Templates["about"])
	c.HTML(200, "application.html", gin.H{
		"wsHost": webUI.Onionhost,
	})
}

func (webUI *WebUI) getContact(c *gin.Context) {
	webUI.Engine.SetHTMLTemplate(webUI.Templates["contact"])
	c.HTML(200, "application.html", gin.H{
		"wsHost": webUI.Onionhost,
	})
}
