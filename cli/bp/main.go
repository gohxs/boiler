package main

//go:generate go get dev.hexasoftware.com/hxs/genversion
//go:generate genversion -package main -out version.go

import (
	"os"

	"github.com/gohxs/boiler/cli/bp/cmd"
)

// main will do the main stuff
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		cmd.RootCmd.Println(err)
		os.Exit(1)
	}
}
