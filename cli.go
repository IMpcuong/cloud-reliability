package main

import (
	cli "github.com/urfave/cli"
)

// Sample commands's flag to run with the specific configuration file (Windows OS only):
// Note: the configuration path can be provided after the flag `-c` or `--config`.
// 	Normal:
// 		.\pdpapp.exe -c .\config\node1\config.json start
// 	Verbose:
//		.\pdpapp.exe --config .\config\node2\config.json start
//	Aliases:
// 		.\pdpapp.exe -c .\config\node2\config.json ims
// 	or
// 		.\pdpapp.exe -c node2 ims (the rest using the same convention)

// newCLIApp create the new CLI application with some custom commands.
func newCLIApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ImChain"
	app.Usage = "Implementation Blockchain in GoLang"

	startServerCLI(app)
	return app
}

// startServerCLI starts the blockchain server and connects to the network.
func startServerCLI(app *cli.App) {
	var cfgPath string
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"ims"},
			Usage:   "start blockchain server",
			Action: func(ctx *cli.Context) error {
				execCmd(ctx, cfgPath)
				if len(ctx.GlobalFlagNames()) > 0 {
					if ctx.String("c") != "" || ctx.String("config") != "" {
						cfgPath = ctx.String("config")
					} else {
						cfgPath = DEFAULT_CFG_PATH
					}
				}
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       DEFAULT_CFG_PATH,
			Usage:       "Load configuration from specific `FILE`",
			Destination: &cfgPath,
		},
	}
}

// execCmd executes the specified commands from the terminal.
func execCmd(ctx *cli.Context, cfgPath string) {
	initNetworkCfg(cfgPath)

	// If `DB_FILE` haven't existed, initialize an empty blockchain.
	// Else, read this file to get the blockchain structure.

	//@@@ FIXME: maybe root cause in here? Unfortunately, the answer is `YES`.
	//@@@ Avoid multi-create/read database file at the same time.
	bc := getLocalBC()
	if bc == nil {
		Info.Printf("Local blockchain database not found. Initialize empty blockchain instead.")
		bc = initBlockChain()
	} else {
		Info.Printf("Import blockchain database from local storage completed!")
	}
	syncNeighborBC(bc)

	if bc == nil || bc.IsEmpty() {
		Info.Printf("Pull failed, no available node for synchronization. Create new blockchain instead.\n")
		bc.AddBlock(newGenesisBlock())
	}

	startBCServer(bc)
	defer bc.DB.Close()
}
