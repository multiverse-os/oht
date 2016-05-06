package oht

import (
	"os"
	"os/signal"
	"syscall"

	"./common"
	"./network"
	"./network/p2p"
	"./network/webui"
)

type OHT struct {
	// Should client name,version and other information be cached here?
	Interface *Interface
	config    *Config
	tor       *network.TorProcess
	p2p       *p2p.Manager
	webUI     *webui.WebUI
	// Channel for shutting down the oht
	//shutdownChan chan bool
	//protocolManager *ProtocolManager -- will this be useful?
	//eventMux *event.TypeMux
}

func NewOHT(torListenPort, torSocksPort, torControlPort, torWebUIPort string) *OHT {
	// Initialize Data Directory
	common.CreatePathUnlessExist("", 0700)
	common.CreatePathUnlessExist("/keys", 0700)
	// Set defaults for torPorts to use if not specified
	config := InitializeConfig(torListenPort, torSocksPort, torControlPort, torWebUIPort)
	// This should be read from the default initialization
	// Should starting tor be a separate function from initialization? Functions to control Tor will be required...
	tor := network.InitializeTor(config.TorListenPort, config.TorSocksPort, config.TorControlPort, config.TorWebUIPort)
	// Initialize WebUI Server
	webUI := webui.InitializeWebUI(tor.WebUIOnionHost, tor.WebUIPort)
	// Initialize & Start P2P Networking
	p2p := p2p.InitializeP2PManager(torListenPort)
	go p2p.Start()
	// Define a clean shutdown process
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tor.Shutdown()
		os.Exit(1)
	}()
	return &OHT{
		Interface: NewInterface(config, tor, webUI, p2p),
		config:    config,
		tor:       tor,
		p2p:       p2p,
		webUI:     webUI,
	}
}
