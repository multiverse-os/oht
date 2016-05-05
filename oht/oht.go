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
	Interface   *Interface
	config      *Config
	tor         *network.TorProcess
	p2p         *p2p.Manager
	webUIServer *webui.Server
	// Channel for shutting down the oht
	//shutdownChan chan bool
	//protocolManager *ProtocolManager -- will this be useful?
	//eventMux *event.TypeMux
}

func NewOHT(torListenPort, torSocksPort, torControlPort, torWebUIPort string) *OHT {
	// Initialize Data Directory
	if !common.FileExist(common.DefaultDataDir()) {
		os.MkdirAll(common.DefaultDataDir(), os.ModePerm)
	}
	// Check if config exists,
	//  if true  - initialize it
	//  if false - load it
	// Set defaults for torPorts to use if not specified
	initializeConfig()
	// This should be read from the default initialization
	config := &Config{
		clientName:         "oht",
		clientMajorVersion: 0,
		clientMinorVersion: 1,
		clientPatchVersion: 0,
	}
	// Should starting tor be a separate function from initialization? Functions to control Tor will be required...
	tor := network.InitializeTor(torListenPort, torSocksPort, torControlPort, torWebUIPort)
	// Start P2P Networking
	p2p := p2p.InitializeP2PManager(torListenPort)
	// WebUI Server
	webUIServer := webui.InitializeServer(tor.WebUIOnionHost, torWebUIPort)
	// Define a clean shutdown process
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tor.Shutdown()
		os.Exit(1)
	}()
	return &OHT{
		Interface:   NewInterface(config, tor, webUIServer, p2p),
		config:      config,
		tor:         tor,
		p2p:         p2p,
		webUIServer: webUIServer,
	}
}
