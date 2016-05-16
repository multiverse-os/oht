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

type TorConfig struct {
	ListenPort  string
	SocksPort   string
	ControlPort string
	WebUIPort   string
}

type Config struct {
	ClientName         string
	ClientMajorVersion string
	ClientMinorVersion string
	ClientPatchVersion string
	MaxPeers           int
	MaxPendingPeers    int
	TorConfig          *TorConfig
	Locale             string
	DevMode            bool
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
		TorConfig: &TorConfig{
			ListenPort:  "9042",
			SocksPort:   "9142",
			ControlPort: "9555",
			WebUIPort:   "8080",
		},
		Locale:       "en",
		DevMode:      false,
		LogFile:      "log.json",
		LogVerbosity: 1,
		Custom:       make(map[string]string),
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
		config.TorConfig.ListenPort = torListenPort
	}
	if torSocksPort != "" {
		config.TorConfig.SocksPort = torSocksPort
	}
	if torControlPort != "" {
		config.TorConfig.ControlPort = torControlPort
	}
	if torWebUIPort != "" {
		config.TorConfig.WebUIPort = torWebUIPort
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

func (c *Config) IPCEndpoint() string {
	// Short circuit if IPC has not been enabled
	if c.IPCPath == "" {
		return ""
	}
	// On windows we can only use plain top-level pipes
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(c.IPCPath, `\\.\pipe\`) {
			return c.IPCPath
		}
		return `\\.\pipe\` + c.IPCPath
	}
	// Resolve names into the data directory full paths otherwise
	if filepath.Base(c.IPCPath) == c.IPCPath {
		if c.DataDir == "" {
			return filepath.Join(os.TempDir(), c.IPCPath)
		}
		return filepath.Join(c.DataDir, c.IPCPath)
	}
	return c.IPCPath
}

// DefaultIPCEndpoint returns the IPC path used by default.
func DefaultIPCEndpoint() string {
	config := &Config{DataDir: common.DefaultDataDir(), IPCPath: common.DefaultIPCSocket}
	return config.IPCEndpoint()
}

// NodeKey retrieves the currently configured private key of the node, checking
// first any manually set key, falling back to the one found in the configured
// data folder. If no key can be found, a new one is generated.
func (c *Config) NodeKey() *ecdsa.PrivateKey {
	// Use any specifically configured key
	if c.PrivateKey != nil {
		return c.PrivateKey
	}
	// Generate ephemeral key if no datadir is being used
	if c.DataDir == "" {
		key, err := crypto.GenerateKey()
		if err != nil {
			glog.Fatalf("Failed to generate ephemeral node key: %v", err)
		}
		return key
	}
	// Fall back to persistent key from the data directory
	keyfile := filepath.Join(c.DataDir, datadirPrivateKey)
	if key, err := crypto.LoadECDSA(keyfile); err == nil {
		return key
	}
	// No persistent key found, generate and store a new one
	key, err := crypto.GenerateKey()
	if err != nil {
		glog.Fatalf("Failed to generate node key: %v", err)
	}
	if err := crypto.SaveECDSA(keyfile, key); err != nil {
		glog.V(logger.Error).Infof("Failed to persist node key: %v", err)
	}
	return key
}

// StaticNodes returns a list of node enode URLs configured as static nodes.
func (c *Config) StaticNodes() []*discover.Node {
	return c.parsePersistentNodes(datadirStaticNodes)
}

// TrusterNodes returns a list of node enode URLs configured as trusted nodes.
func (c *Config) TrusterNodes() []*discover.Node {
	return c.parsePersistentNodes(datadirTrustedNodes)
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
