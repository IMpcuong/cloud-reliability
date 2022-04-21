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

// NewCLIApp create the new CLI application with some custom commands.
func NewCLIApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ImChain"
	app.Usage = "Implementation Blockchain in GoLang"

	StartServerCLI(app)
	return app
}

// StartServerCLI starts the blockchain server and connects to the network.
func StartServerCLI(app *cli.App) {
	var cfgPath string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       DEFAULT_CFG_PATH,
			Usage:       "Load configuration from specific `FILE`",
			Destination: &cfgPath,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"ims"},
			Usage:   "start blockchain server",
			Action: func(ctx *cli.Context) error {
				ExecCmd(ctx, cfgPath)
				return nil
			},
		},
	}
}

// ExecCmd executes the specified commands from the terminal.
func ExecCmd(ctx *cli.Context, cfgPath string) {
	InitNetworkCfg(cfgPath)

	bc := PullNeighborBC()
	if bc == nil || bc.IsEmpty() {
		Info.Printf("Pull failed. Create new blockchain instead.\n")
		bc = InitBlockChain()
	}
	StartBCServer(bc)
}
