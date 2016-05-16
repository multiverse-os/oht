package oht

import (
	"os"
	"os/signal"
	"syscall"

	"./common"
	"./network"
	"./network/webui"
)

type OHT struct {
	Interface *Interface
	config    *Config
	tor       *network.TorProcess
	p2p       *network.Manager
	webUI     *webui.WebUI
	Shutdown  chan os.Signal
}

func (oht *OHT) cleanShutdown(c chan os.Signal) {
	signal.Notify(oht.Shutdown, os.Interrupt, syscall.SIGTERM)
	<-c
	oht.Stop()
}

func NewOHT(torListenPort, torSocksPort, torControlPort, torWebUIPort string) (oht *OHT) {
	common.CreatePathUnlessExist("", 0700)
	common.CreatePathUnlessExist("/keys", 0700)
	config := InitializeConfig(torListenPort, torSocksPort, torControlPort, torWebUIPort)
	tor := network.InitializeTor(config)
	p2p := p2p.InitializeP2PManager(config)
	webUI := webui.InitializeWebUI(tor.WebUIOnionHost, config.TorWebUIPort)
	oht = &OHT{
		Interface: NewInterface(config, tor, webUI, p2p),
		config:    config,
		tor:       tor,
		p2p:       p2p,
		webUI:     webUI,
		Shutdown:  make(chan os.Signal, 1),
	}
	go oht.cleanShutdown(oht.Shutdown)
	return oht
}

func (oht *OHT) Start() bool {
	oht.tor.Start()
	go oht.network.Start()
	return true
}

func (oht *OHT) Stop() bool {
	oht.webUI.Server.Stop()
	// Stop p2p
	oht.network.Stop()
	os.Exit(1)
	return true
}
