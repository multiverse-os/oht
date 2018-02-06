package oht

import (
	"os"
	"os/signal"
	"syscall"

	"lib/oht/core/common"
	"lib/oht/core/network"
	"lib/oht/core/network/webui"
)

// Need to scrap the current struct with pointers for a nested interface system
// Networking, P2P, Tor libraries should be emitting events that this file can listen in on
// then handle those events in a switch case surrounded by a for loop.
type OHT struct {
	// None of these points are concurrent safe. Where possible they should be converted to values
	Interface *Interface
	config    *config.Config
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
	config := InitializeConfig(torListenPort, torSocksPort, torControlPort, torWebUIPort)
	common.CreatePathUnlessExist(config.DataDirectory+"", 0700)
	common.CreatePathUnlessExist(config.DataDirectory+"keys", 0700)
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
	go oht.p2p.Start()
	return true
}

func (oht *OHT) Stop() bool {
	oht.webUI.Server.Stop()
	oht.tor.Stop(false)
	oht.p2p.Stop()
	os.Exit(1)
	return true
}
