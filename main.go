package main

import (
	"github.com/antonbiluta/blockchain-golang/cli"
	"os"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
