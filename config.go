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
	DEFAULT_NW_NODES = 3
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
func getNetworkCfg() Config {
	return nwConfig
}

// initNetworkCfg initializes the network configurations from the config file
// from the given source path.
func initNetworkCfg(cfgPathCLI string) {
	var cfgPath string
	if cfgPathCLI != "" {
		cfgPath = cfgPathCLI
	} else {
		cfgPath = DEFAULT_CFG_PATH
	}
	nwConfig = importNetworkCfg(cfgPath)
}

// importNetworkCfg reads the configuration from file in given `path` (flag value)
// and returns the network configuration.
func importNetworkCfg(path string) Config {
	var cfgPath string
	// Walk through the default config directory and returns all the sub-directories.
	dirs, err := walkCfgDir("")
	if err != nil {
		Error.Println(err.Error())
	}
	// Read all the sub-directories and returns relative paths
	// from available config files for each nodes.
	paths, err := readPaths(dirs)
	if err != nil {
		Error.Println(err.Error())
	}
	curNodes := len(dirs)
	if curNodes < DEFAULT_NW_NODES {
		curNodes = DEFAULT_NW_NODES
	}
	for node := 1; node <= curNodes; node++ {
		// Checking if the flag value equal to the string formatter (n1/2/3) or not.
		flag := fmt.Sprintf("%s%d", "n", node)
		if path == flag {
			// NOTE: config.json maybe change the position later,
			// contemporary take the second place in total 3 config files in each node.
			cfgPath = paths[node*3-2]
		}
	}
	cfgFile, err := ioutil.ReadFile(cfgPath)
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

// walkCfgDir walks through the directory tree structure and returns all the sub-directories.
func walkCfgDir(cfgDir string) ([]string, error) {
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
	if contains(dirs, DEFAULT_CFG_DIR) {
		dirs = remove(dirs, DEFAULT_CFG_DIR)
	}
	if err != nil {
		Error.Println(err)
	}
	return unique(dirs), err
}

// readPaths reads multiple directories as an argument
// and return the relative path of all config files.
func readPaths(dirs []string) ([]string, error) {
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
				filePaths, _ = readPaths([]string{dir + "/" + file.Name()})
			}
		}
	}
	return filePaths, nil
}
