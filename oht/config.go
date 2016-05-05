package oht

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"./common"
)

var (
	// This needs to be implenetd, a standard model for a node and a method of both parsing and writing nodes
	staticNodes  = "static-nodes.json"  // Path within <datadir> to search for the static node list
	trustedNodes = "trusted-nodes.json" // Path within <datadir> to search for the trusted node list
)

type Config struct {
	ClientName         string
	ClientMajorVersion string
	ClientMinorVersion string
	ClientPatchVersion string

	MaxPeers        int
	MaxPendingPeers int

	TorListenPort  string
	TorSocksPort   string
	TorControlPort string
	TorWebUIPort   string

	Locale string

	DevMode bool

	DataDir     string
	GenesisFile string `json:",omitempty"`

	LogFile      string
	LogVerbosity int

	Custom map[string]string
	// NewDB is used to create databases.
	// If nil, the default is to create boltdb databases on disk.
	// -- Setup the boltDB here
}

func InitializeConfig(torListenPort, torSocksPort, torControlPort, torWebUIPort string) (config *Config) {
	config = &Config{
		ClientName:         "oht",
		ClientMajorVersion: "0",
		ClientMinorVersion: "1",
		ClientPatchVersion: "0",
		MaxPeers:           8,
		MaxPendingPeers:    8,
		TorListenPort:      "9042",
		TorSocksPort:       "9142",
		TorControlPort:     "9555",
		TorWebUIPort:       "8080",
		Locale:             "en",
		DevMode:            false,
		DataDir:            common.DefaultDataDir(),
		LogFile:            (common.DefaultDataDir() + "/log.json"),
		LogVerbosity:       1,
		Custom:             make(map[string]string),
	}
	if _, err := ioutil.ReadFile(common.AbsolutePath(common.DefaultDataDir(), "config.json")); err != nil {
		jsonFile, err := json.Marshal(config)
		if err = ioutil.WriteFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"), jsonFile, 0644); err != nil {
			log.Fatal(err)
		}
	}
	jsonFile, err := ioutil.ReadFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"))
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

func (config *Config) saveConfiguration() bool {
	jsonFile, err := json.Marshal(config)
	if err = ioutil.WriteFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"), jsonFile, 0644); err != nil {
		return false
	}
	return true
}

func (config *Config) clientInfo() string {
	return common.CompileClientInfo(config.ClientName, config.clientVersion())
}

func (config *Config) clientVersion() string {
	return config.ClientMajorVersion + "." + config.ClientMinorVersion + "." + config.ClientPatchVersion
}

func (config *Config) setConfigOption(key, value string) bool {
	// Validate length of value
	if key == "ClientName" {
		config.ClientName = value
	} else if key == "ClientMajorVersion" {
		config.ClientMajorVersion = value
	} else if key == "ClientMinorVersion" {
		config.ClientMinorVersion = value
	} else if key == "ClientPatchVersion" {
		config.ClientPatchVersion = value
	} else if key == "MaxPeers" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return false
		} else {
			config.MaxPeers = i
		}
	} else if key == "MaxPendingPeers" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return false
		} else {
			config.MaxPendingPeers = i
		}
	} else if key == "TorListenPort" {
		config.TorListenPort = value
	} else if key == "TorSocksPort" {
		config.TorSocksPort = value
	} else if key == "TorControlPort" {
		config.TorControlPort = value
	} else if key == "TorWebUIPort" {
		config.TorWebUIPort = value
	} else if key == "TorListenPort" {
		config.TorListenPort = value
	} else if key == "Locale" {
		config.Locale = value
	} else if key == "DevMode" {
		b, err := strconv.ParseBool("true")
		if err != nil {
			return false
		} else {
			config.DevMode = b
		}
	} else if key == "DataDir" {
		config.DataDir = value
	} else if key == "LogFile" {
		config.LogFile = value
	} else if key == "LogVerbosity" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return false
		} else {
			config.LogVerbosity = i
		}
	} else {
		config.Custom[key] = value
	}
	return true
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
