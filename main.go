package main

import (
	"os"

	"github.com/jcsuscriptor/elmahbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
