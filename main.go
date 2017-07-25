package main

//go:generate go get dev.hexasoftware.com/hxs/genversion
//go:generate genversion -package main -out version.go

import (
	"fmt"
	"os"

	"github.com/gohxs/boiler/internal/core"
)

func main() {
	app := core.NewApp(os.Stdin)
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
