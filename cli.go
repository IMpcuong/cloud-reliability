package main

import (
	"fmt"
	"os"

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

	Sample command that will create wallet and store its contents to 'config.json' file:
	Normal:
		.\pdpapp.exe --wa node1 create-wallet
	Verbose:
		.\pdpapp.exe --wallet-addr node1 create-wallet
	Alias:
		.\pdpapp.exe --wa node1 cw

	Sample command of creating new transaction: (eg: send from node2 -> node1)
	Step 1: .\pdpapp.exe -c node2 -n node2 ims
	Step 2: .\pdpapp.exe crtx -c node2 -n node2 --to localhost:3331 -v 1 -f test.txt
	Step 3: Checking the `test.txt` file for more details.
*/

// newCLIApp create the new CLI application with some custom commands.
func newCLIApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ImChain"
	app.Usage = "Implementation Blockchain in GoLang"

	// NOTE: global commands and flags.
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}

	startServerCLI(app)
	createWalletCLI(app)
	createTransactionCLI(app)

	return app
}

func createWalletCLI(app *cli.App) {
	var cfgPath string

	app.Commands = append(app.Commands, []cli.Command{
		{
			Name:    "create-wallet",
			Aliases: []string{"cw"},
			Usage:   "create new storable wallet address",
			Action: func(ctx *cli.Context) error {
				execCreateWallet(ctx, cfgPath)
				return nil
			},
		},
	}...)
	app.Flags = append(app.Flags, []cli.Flag{
		cli.StringFlag{
			Name:        "wallet-addr, wa",
			Value:       DEFAULT_CFG_PATH,
			Usage:       "Export Wallet's configuration to specific `FILE`",
			Destination: &cfgPath,
		},
	}...)
}

// startServerCLI starts the blockchain server and connects to the network.
func startServerCLI(app *cli.App) {
	var cfgPath, nodeDb string

	app.Commands = append(app.Commands, []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"ims"},
			Usage:   "start blockchain server",
			Action: func(ctx *cli.Context) error {
				execStartServer(ctx, cfgPath, nodeDb)
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
	}...)
	app.Flags = append(app.Flags, []cli.Flag{
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
	}...)
}

func createTransactionCLI(app *cli.App) {
	var cfgPath, nodeDb, toAddr, exportFile string
	var totalVal int

	app.Commands = append(app.Commands, []cli.Command{
		{
			Name:    "create-tx",
			Aliases: []string{"crtx"},
			Usage:   "crtx -c {cfgPath} -n {node} -to {address} -v {value} -f {exportFile}",
			Action: func(ctx *cli.Context) error {
				execCreateTx(ctx, totalVal, cfgPath, nodeDb, toAddr, exportFile)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "c",
					// cfgPath[0]
					Destination: &cfgPath,
				},
				cli.StringFlag{
					Name: "n",
					// cfgPath[1]
					Destination: &nodeDb,
				},
				cli.IntFlag{
					Name:        "v",
					Destination: &totalVal,
				},
				cli.StringFlag{
					Name: "to",
					// cfgPath[2]
					Destination: &toAddr,
				},
				cli.StringFlag{
					Name: "f",
					// cfgPath[3]
					Destination: &exportFile,
				},
			},
		},
	}...)
}

// execStartServer executes the specified commands from the terminal.
func execStartServer(ctx *cli.Context, cfgPath ...string) {
	// `cfg[0]` = path to the configuration file.
	// `cfg[1]` = path to the database storage file.
	initNwCfg(cfgPath[0])

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
		firstTx := []Transaction{}
		bc.AddBlock(newGenesisBlock(firstTx))
	}

	startBCServer(bc)
	closeDB(bc)
}

// execCreateWallet creates new a `Wallet` instance.
func execCreateWallet(ctx *cli.Context, cfgPath string) {
	config := initNwCfg(cfgPath)
	wallet := newWallet()
	config.WJson = *wallet.ToJson()
	config.ExportNetworkCfg(cfgPath)

	fmt.Printf("New wallet is created successfully! Wallet is exported to : * %s *\n", cfgPath)
	fmt.Printf("%s\n", config.WJson)
}

// @@@ FIXME: to be more cleaner!
func execCreateTx(ctx *cli.Context, val int, cfgPath ...string) {
	// `cfg[0]` = path to the configuration file.
	// `cfg[1]` = path to the database storage file.
	// `cfg[2]` = the node's port address.
	// `cfg[3]` = the path to the output file.
	Info.Printf("Execute transaction: send %d coins to address %s", val, cfgPath[2])
	initNwCfg(cfgPath[0])
	wallet := getWallet()
	bc := getLocalBC(cfgPath[1])
	if bc == nil {
		Error.Print("Local blockchain not found. Need one existed first!")
		os.Exit(1)
	}

	tx := bc.NewTx(wallet, cfgPath[2], val)
	msgReq := createMsgReqAddTx(tx)
	msgReq.Export(cfgPath[3])
	closeDB(bc)
}
