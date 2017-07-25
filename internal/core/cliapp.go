package core

import (
	"io"

	"github.com/urfave/cli"
)

// NewApp return cli app
func NewApp(stdin io.Reader) *cli.App {

	app := cli.NewApp()
	app.Name = "boiler"
	app.Usage = "Generic boilerplate app"
	app.Author = "Luis Figueiredo"
	app.Email = "luisf@hexasoftware.com"
	app.Version = "0.0.1"

	app.Metadata = map[string]interface{}{
		"stdin": stdin,
	}

	app.Commands = []cli.Command{
		commandInit(),
		commandGenerate(),
	}

	return app

}
