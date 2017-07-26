package main

//go:generate go get dev.hexasoftware.com/hxs/genversion
//go:generate genversion -package main -out version.go

import (
	"fmt"
	"log"
	"os"

	"github.com/gohxs/boiler/internal/core"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	app := core.NewApp(os.Stdin)
	app.Version = Version
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
