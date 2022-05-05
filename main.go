package main

import (
	"io/ioutil"
	"os"
)

func main() {
	// UNUSED
	_ = deserializeChain(new(Blockchain).Serialize())
	// UNUSED

	GenerateLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	cliApp := newCLIApp()
	cliApp.Run(os.Args)
}
