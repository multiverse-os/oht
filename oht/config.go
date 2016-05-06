package oht

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"./common"
)

var (
	configOptions = common.Attributes(&Config{})
)

type Config struct {
	ClientName         string
	ClientMajorVersion string
	ClientMinorVersion string
	ClientPatchVersion string
	MaxPeers           int
	MaxPendingPeers    int
	TorListenPort      string
	TorSocksPort       string
	TorControlPort     string
	TorWebUIPort       string
	Locale             string
	DevMode            bool
	KnownNodesFile     string
	TrustedNodesFile   string
	DataDir            string
	GenesisFile        string `json:",omitempty"`
	LogFile            string
	LogVerbosity       int
	Custom             map[string]string
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
		KnownNodesFile:     "known-nodes.json",
		TrustedNodesFile:   "trusted-nodes.json",
		LogFile:            "log.json",
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

func (config *Config) zeroValue(key string) bool {
	if len(key) >= 255 || key == "Custom" {
		return false
	}
	cvalue := reflect.ValueOf(config)
	copyValue := cvalue.Elem()
	typeField := copyValue.FieldByName(key).Type()
	zero := reflect.Zero(typeField)
	copyValue.FieldByName(key).Set(zero)
	return true
}

func (config *Config) setValue(key, value string) bool {
	if len(key) >= 255 || key == "Custom" {
		return false
	}
	cvalue := reflect.ValueOf(config)
	copyValue := cvalue.Elem()
	field := copyValue.FieldByName(key)
	typeField := field.Type()
	if typeField == reflect.TypeOf(value) {
		copyValue.FieldByName(key).SetString(value)
	} else if typeField == reflect.TypeOf(0) {
		i, err := strconv.Atoi(value)
		if err != nil {
			return false
		} else {
			copyValue.FieldByName(key).SetInt(int64(i))
		}
	} else if typeField == reflect.TypeOf(true) {
		b, err := strconv.ParseBool("true")
		if err != nil {
			return false
		} else {
			copyValue.FieldByName(key).SetBool(b)
		}
	}
	return true
}

func (config *Config) unsetConfigOption(key string) bool {
	if len(key) >= 255 || key == "Custom" {
		return false
	}
	if common.ItemInSlice(key, configOptions) {
		config.zeroValue(key)
	} else {
		delete(config.Custom, key)
	}
	return true
}

func (config *Config) setConfigOption(key, value string) bool {
	if (len(key) >= 255 && len(value) >= 255) || key == "Custom" {
		return false
	}
	if common.ItemInSlice(key, configOptions) {
		return config.setValue(key, value)
	} else {
		config.Custom[key] = value
		return true
	}
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
