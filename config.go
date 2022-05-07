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
var (
	nwConfig *Config
)

// Required configurations for the network.
type Config struct {
	Network Network    `json:"network"` // Network configurations.
	WJson   WalletJson `json:"wallet"`  // Address identification properties.
}

// Utility functions start from here.

// ExportNetworkCfg export the network configuration to the given file path.
func (cfg *Config) ExportNetworkCfg(filePath string) {
	prettyMarshal, e := json.MarshalIndent(cfg, "", "  ")
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}

	e = ioutil.WriteFile(filePath, prettyMarshal, 0644)
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}
}

// Get default network configurations.
func getNetworkCfg() *Config {
	return nwConfig
}

// initNetworkCfg initializes the network configurations from the config file
// with the given source path.
func initNetworkCfg(cfgPathCLI string) *Config {
	var cfgPath string
	if cfgPathCLI != "" {
		cfgPath = cfgPathCLI
	} else {
		cfgPath = DEFAULT_CFG_PATH
	}
	nwConfig = importNetworkCfg(cfgPath)
	return nwConfig
}

// importNetworkCfg reads the configuration from file in given `path` (flag value)
// and returns the network configuration.
func importNetworkCfg(path string) *Config {
	if (strings.Compare(path, DEFAULT_CFG_PATH)) == 0 {
		return getCfgData(path)
	}

	var cfgPath string
	// Walk through the default config directory and returns all the sub-directories.
	dirs, err := walkCfgDir("")
	if err != nil {
		Error.Println(err.Error())
	}

	// Read all the sub-directories and returns relative paths
	// from available config files for each nodes.
	for iNode, dir := range dirs {
		paths, err := readPaths(dir)
		if err != nil {
			Error.Println(err.Error())
		}

		for _, filePath := range paths {
			// Checking if the flag value equal to the string formatter (node1/2/3) or not.
			flag := fmt.Sprintf("%s%d", "node", iNode+1)
			if path == flag && strings.Contains(filePath, "config.json") {
				cfgPath = filePath
			}
		}
	}
	return getCfgData(cfgPath)
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

// readPaths reads given directory as an argument
// and return the relative path of all config files.
func readPaths(dir string) ([]string, error) {
	var filePaths []string

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return filePaths, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", dir, file.Name())
			filePaths = append(filePaths, path)
		} else {
			filePaths, _ = readPaths(dir + "/" + file.Name())
		}
	}

	return filePaths, nil
}

// readFile reads the file contents from the given path.
func readFile(path string) ([]byte, error) {
	cfgFile, err := ioutil.ReadFile(path)
	if err != nil {
		Error.Println(err.Error())
		os.Exit(1)
	}
	return cfgFile, nil
}

// getCfgData returns the configuration data corresponding to the given path.
func getCfgData(path string) *Config {
	cfgFile, _ := readFile(path)
	cfg := Config{}
	err := json.Unmarshal(cfgFile, &cfg)
	if err != nil {
		Error.Println(err.Error())
		os.Exit(1)
	}
	return &cfg
}
