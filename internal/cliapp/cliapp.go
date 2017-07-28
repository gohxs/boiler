package cliapp

import (
	"io"

	"github.com/gohxs/boiler/internal/core"
	"github.com/urfave/cli"
)

var (
	gCore, _ = core.FromCurDir()
)

// NewApp return cli app
func NewApp(stdin io.Reader) *cli.App {

	// Load root/app core context thing here
	// Plus data for app if inside a boiler proj

	app := cli.NewApp()
	app.Name = "boiler"
	app.Usage = "Generic boilerplate app"
	app.Author = "Luis Figueiredo"
	app.Email = "luisf@hexasoftware.com"

	app.Metadata = map[string]interface{}{
		"stdin": stdin,
	}

	app.Commands = []cli.Command{
		commandNew(),
		commandCreate(),
		commandCreatebak(),
	}

	return app

}
