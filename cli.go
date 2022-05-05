package main

import (
	cli "github.com/urfave/cli"
)

/*
	Sample commands's flag to run with the specific configuration and database storage file:
	Note: the configuration path can be provided after the flag `-c` or `--config`, `.exe`
	      to make the binaries file executable in Windows environment.
	Normal:
 		.\pdpapp.exe -c node1 -n node1 start
 	Verbose:
		.\pdpapp.exe --config node1 --node node1 start
	Aliases:
 		.\pdpapp.exe -c node2 -n node2 ims
*/

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
	var cfgPath, nodeDb string

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"ims"},
			Usage:   "start blockchain server",
			Action: func(ctx *cli.Context) error {
				execCmd(ctx, cfgPath, nodeDb)
				if len(ctx.GlobalFlagNames()) > 0 {
					if ctx.String("c") != "" {
						cfgPath = ctx.String("c")
					} else if ctx.String("config") != "" {
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
		cli.StringFlag{
			Name:        "node, n",
			Value:       "",
			Usage:       "Load database storage from specified `NODE`",
			Destination: &nodeDb,
		},
	}
}

// execCmd executes the specified commands from the terminal.
func execCmd(ctx *cli.Context, cfgPath ...string) {
	// `cfg[0]` = path to the configuration file.
	// `cfg[1]` = path to the database storage file.
	initNetworkCfg(cfgPath[0])

	// If `DB_FILE` haven't existed, initialize an empty blockchain.
	// Else, read this file to get the blockchain structure.
	bc := getLocalBC(cfgPath[1])
	if bc == nil {
		Info.Printf("Local blockchain database not found. Initialize empty blockchain instead.")
		bc = initBlockChain(cfgPath[1])
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
