package oht

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"./common"
	"./network"
	"./network/p2p"
	"./network/webui"
)

type OHT struct {
	Interface *Interface
	config    *Config
	tor       *network.TorProcess
	p2p       *p2p.Manager
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
	log.Println("configed")
	tor := network.InitializeTor(config.TorListenPort, config.TorSocksPort, config.TorControlPort, config.TorWebUIPort)
	log.Println("torred")
	p2p := p2p.InitializeP2PManager(config.TorListenPort)
	log.Println("p2ped")
	webUI := webui.InitializeWebUI(tor.WebUIOnionHost, config.TorWebUIPort)
	log.Println("webuid")
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
	go oht.p2p.Start()
	return true
}

func (oht *OHT) Stop() bool {
	// Stop webUI
	oht.webUI.Server.Stop()
	// Stop p2p
	// Stop tor
	oht.tor.Stop()
	os.Exit(1)
	return true
}
