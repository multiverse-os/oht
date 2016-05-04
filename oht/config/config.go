package config

import (
	"io/ioutil"
	"log"

	"../common"
)

type Config struct {
	// Load this struct from the config.json file
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
