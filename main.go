package main

import (
	"io/ioutil"
	"os"
)

func main() {
	GenerateLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	cliApp := newCLIApp()
	cliApp.Run(os.Args)
}
