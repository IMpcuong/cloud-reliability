package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_CFG_DIR  = "config"
	DEFAULT_CFG_PATH = "config/config.json"
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

func WalkCfgDir(cfgDir string) ([]string, error) {
	if strings.Compare(cfgDir, "") == 0 {
		cfgDir = DEFAULT_CFG_DIR
	}
	var dirs []string
	err := filepath.Walk(cfgDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, strings.ReplaceAll(path, `\`, `/`))
		}
		return nil
	})
	if err != nil {
		Error.Println(err)
	}
	return dirs, err
}

func ReadPaths(dirs []string) ([]string, error) {
	// if strings.Compare(strings.Join(dirs, ""), "") == 0 {
	// 	dirs = []string{DEFAULT_CFG_DIR}
	// }
	var filePaths []string
	for _, dir := range dirs {
		fileInfo, err := ioutil.ReadDir(dir)
		if err != nil {
			return filePaths, err
		}
		for _, file := range fileInfo {
			if !file.IsDir() {
				path := fmt.Sprintf("%s/%s", dir, file.Name())
				filePaths = append(filePaths, path)
			} else {
				filePaths, err = ReadPaths([]string{dir + "/" + file.Name()})
			}
		}
	}
	return filePaths, nil
}
