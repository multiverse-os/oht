package oht

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"./common"
)

var (
	staticNodes  = "static-nodes.json"  // Path within <datadir> to search for the static node list
	trustedNodes = "trusted-nodes.json" // Path within <datadir> to search for the trusted node list
)

type Config struct {
	ClientName         string
	ClientMajorVersion int
	ClientMinorVersion int
	ClientPatchVersion int

	TorListenPort  string
	TorSocksPort   string
	TorControlPort string
	TorWebUIPort   string

	ProtocolVersion int `json:",omitempty"`

	DevMode bool `json:",omitempty"`
	TestNet bool `json:",omitempty"`

	NetworkId   int    `json:",omitempty"`
	GenesisFile string `json:",omitempty"`

	DatabaseCache int `json:",omitempty"`

	DataDir   string `json:",omitempty"`
	LogFile   string `json:",omitempty"`
	Verbosity int    `json:",omitempty"`
	ExtraData []byte `json:",omitempty"`

	MaxPeers        int  `json:",omitempty"`
	MaxPendingPeers int  `json:",omitempty"`
	Discovery       bool `json:",omitempty"`
	// NewDB is used to create databases.
	// If nil, the default is to create boltdb databases on disk.
	// -- Setup the boltDB here
}

func InitializeConfig(torListenPort, torSocksPort, torControlPort, torWebUIPort string) (config *Config) {
	config = &Config{
		ClientName:         "oht",
		ClientMajorVersion: 0,
		ClientMinorVersion: 1,
		ClientPatchVersion: 0,
		TorListenPort:      "9042",
		TorSocksPort:       "9142",
		TorControlPort:     "9555",
		TorWebUIPort:       "8080",
	}
	jsonFile, err := ioutil.ReadFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"))
	if err != nil {
		jsonFile, err := json.Marshal(config)
		if err = ioutil.WriteFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"), jsonFile, 0644); err != nil {
			log.Fatal(err)
		}
	}
	err = json.Unmarshal(jsonFile, config)
	if err != nil {
		log.Fatal(err)
	}
	if torListenPort != "" {
		config.TorListenPort = torListenPort
	}
	if torSocksPort != "" {
		config.TorSocksPort = torSocksPort
	}
	if torControlPort != "" {
		config.TorControlPort = torControlPort
	}
	if torWebUIPort != "" {
		config.TorWebUIPort = torWebUIPort
	}
	return config
}

func (config *Config) clientInfo() string {
	return common.CompileClientInfo(config.ClientName, config.clientVersion())
}

func (config *Config) clientVersion() string {
	return fmt.Sprintf("%d.%d.%d", config.ClientMajorVersion, config.ClientMinorVersion, config.ClientPatchVersion)
}

// parseNodes parses a list of discovery node URLs loaded from a .json file.
//func (cfg *Config) parseNodes(file string) []*discover.Node {
//	// Short circuit if no node config is present
//	path := filepath.Join(cfg.DataDir, file)
//	if _, err := os.Stat(path); err != nil {
//		return nil
//	}
//	// Load the nodes from the config file
//	blob, err := ioutil.ReadFile(path)
//	if err != nil {
//		log.Println("Failed to access nodes: %v", err)
//		return nil
//	}
//	nodelist := []string{}
//	if err := json.Unmarshal(blob, &nodelist); err != nil {
//		log.Println("Failed to load nodes: %v", err)
//		return nil
//	}
//	// Interpret the list as a discovery node array
//	var nodes []*discover.Node
//	for _, url := range nodelist {
//		if url == "" {
//			continue
//		}
//		node, err := discover.ParseNode(url)
//		if err != nil {
//			log.Println("Node URL %s: %v\n", url, err)
//			continue
//		}
//		nodes = append(nodes, node)
//	}
//	return nodes
//}
//
//func (cfg *Config) nodeKey() (*ecdsa.PrivateKey, error) {
//	// use explicit key from command line args if set
//	if cfg.NodeKey != nil {
//		return cfg.NodeKey, nil
//	}
//	// use persistent key if present
//	keyfile := filepath.Join(cfg.DataDir, "nodekey")
//	key, err := crypto.LoadECDSA(keyfile)
//	if err == nil {
//		return key, nil
//	}
//	// no persistent key, generate and store a new one
//	if key, err = crypto.GenerateKey(); err != nil {
//		return nil, fmt.Errorf("could not generate server key: %v", err)
//	}
//	if err := crypto.SaveECDSA(keyfile, key); err != nil {
//		log.PrintLn("could not persist nodekey: ", err)
//	}
//	return key, nil
//}
