package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	DEFAULT_CFG_PATH = "./config/config.json"
)

// Default variable of network configuration.
var nwConfig Config

// Required configurations for the network.
type Config struct {
	Network Network `json:"network"`
}

// Utility functions start from here.

// Get default network configurations.
func GetNetworkCfg() Config {
	return nwConfig
}

// InitNetworkCfg initializes the network configurations from the config file
// from the given source path.
func InitNetworkCfg(cfgPathCLI string) {
	var cfgPath string
	if cfgPathCLI != "" {
		cfgPath = cfgPathCLI
	} else {
		cfgPath = DEFAULT_CFG_PATH
	}
	nwConfig = ImportNetworkCfg(cfgPath)
}

// ImportNetworkCfg reads the configuration from file in given `path`
// and returns the network configuration.
func ImportNetworkCfg(path string) Config {
	cfgFile, err := ioutil.ReadFile(path)
	if err != nil {
		Error.Println(err.Error())
		os.Exit(1)
	}

	cfg := Config{}
	err = json.Unmarshal(cfgFile, &cfg)
	if err != nil {
		Error.Println(err.Error())
		os.Exit(1)
	}
	return cfg
}
