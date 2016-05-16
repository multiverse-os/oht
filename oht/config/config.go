package config

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"../common"
	"../crypto"
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
	DataDirectory      string
	IPCName            string
	IPCPath            string
	GenesisFile        string `json:",omitempty"`
	PrivateKeyFile     string
	PrivateKey         *ecdsa.PrivateKey
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
		Locale:         "en",
		IPCName:        "oht",
		DevMode:        false,
		LogFile:        "log.json",
		DataDirectory:  common.DefaultDataDir("oht"),
		PrivateKeyFile: (common.DefaultDataDir("oht") + "node_key"),
		LogVerbosity:   1,
		Custom:         make(map[string]string),
	}
	if _, err := ioutil.ReadFile(common.AbsolutePath(config.DataDirectory, "config.json")); err != nil {
		jsonFile, err := json.Marshal(config)
		if err = ioutil.WriteFile(common.AbsolutePath(config.DataDirectory, "config.json"), jsonFile, 0644); err != nil {
			log.Fatal(err)
		}
	}
	jsonFile, err := ioutil.ReadFile(common.AbsolutePath(config.DataDirectory, "config.json"))
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
	if err = ioutil.WriteFile(common.AbsolutePath(config.DataDirectory, "config.json"), jsonFile, 0644); err != nil {
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

func (config *Config) IPCEndpoint() string {
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(config.IPCPath, `\\.\pipe\`) {
			return config.IPCPath
		}
		return `\\.\pipe\` + config.IPCPath
	}
	if filepath.Base(config.IPCPath) == config.IPCPath {
		if config.DataDirectory == "" {
			return filepath.Join(os.TempDir(), config.IPCPath)
		}
		return filepath.Join(config.DataDirectory, config.IPCPath)
	}
	return config.IPCPath
}

func (config *Config) NodeKey() *ecdsa.PrivateKey {
	if config.PrivateKey != nil {
		if common.FileExist(config.PrivateKeyFile) {
			if key, err := crypto.LoadECDSA(config.PrivateKeyFile); err == nil {
				config.PrivateKey = key
			}
		} else {
			if key, err := crypto.GenerateKey(); err == nil {
				if err := crypto.SaveECDSA(config.PrivateKeyFile, key); err != nil {
					log.Println("Config: Failed to save node key.")
				}
				config.PrivateKey = key
			}
		}
	}
	return config.PrivateKey
}

func (config *Config) StaticNodes() []*discover.Node {
	return config.parsePersistentNodes(datadirStaticNodes)
}

func (config *Config) TrusterNodes() []*discover.Node {
	return config.parsePersistentNodes(datadirTrustedNodes)
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
