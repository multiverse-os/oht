package config

import (
	"crypto/ecdsa"
	"io/ioutil"
	"log"
	"regexp"

	"../../accounts"
	"../common"
)

var (
	datadirInUseErrnos = map[uint]bool{11: true, 32: true, 35: true}
	portInUseErrRE     = regexp.MustCompile("address already in use")

	staticNodes  = "static-nodes.json"  // Path within <datadir> to search for the static node list
	trustedNodes = "trusted-nodes.json" // Path within <datadir> to search for the trusted node list
)

type Config struct {
	clientVersion string
	netVersionId  int
	// Load this struct from the config.json file
	DevMode bool
	TestNet bool

	Name        string
	NetworkId   int
	GenesisFile string

	DatabaseCache int

	DataDir   string
	LogFile   string
	Verbosity int
	ExtraData []byte

	MaxPeers        int
	MaxPendingPeers int
	Discovery       bool
	// Need specific ports - ... this is a string, can I just use strings?
	//Port            string

	// This key is used to identify the node on the network.
	// If nil, an ephemeral key is used.
	NodeKey *ecdsa.PrivateKey

	Etherbase      common.Address
	AccountManager *accounts.Manager

	// NewDB is used to create databases.
	// If nil, the default is to create boltdb databases on disk.
	// -- Setup the boltDB here
}

func InitializeConfig() {
	if _, err := ioutil.ReadFile(common.AbsolutePath(common.DefaultDataDir(), "config.json")); err != nil {
		str := "{}"
		if err = ioutil.WriteFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"), []byte(str), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func DisplayConfig() (configData []byte, err error) {
	configData, err = ioutil.ReadFile(common.AbsolutePath(common.DefaultDataDir(), "config.json"))
	return configData, err
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
