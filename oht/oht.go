package oht

import (
	"os"

	"./../accounts"
	"./../contacts"
	"./common"
)

type OHT struct {
	// Should client name,version and other information be cached here?
	Interface *Interface
	config    *Config
	// Channel for shutting down the oht
	shutdownChan chan bool
	// This should be in its own sub-module
	accountManager *accounts.Manager
	//protocolManager *ProtocolManager -- will this be useful?
	//eventMux *event.TypeMux
}

func NewOHT() *OHT {
	// Initialize Data Directory
	if !common.FileExist(common.DefaultDataDir()) {
		os.MkdirAll(common.DefaultDataDir(), os.ModePerm)
	}
	// Check if config exists,
	//  if true  - initialize it
	//  if false - load it
	InitializeConfig()
	contacts.InitializeContacts()
	config := &Config{
		clientName:         "oth",
		clientMajorVersion: 0,
		clientMinorVersion: 1,
		clientPatchVersion: 0,
	}
	accountsManager := &accounts.Manager{}
	return &OHT{
		Interface:      NewInterface(config, accountsManager),
		config:         config,
		accountManager: accountsManager,
	}
}
