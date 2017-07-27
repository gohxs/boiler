package cliapp

import (
	"log"

	"github.com/urfave/cli"
)

func commandCreate() cli.Command {
	ret := cli.Command{
		Name:  "create",
		Usage: "Create a thing from template",
	}

	subCommands := []cli.Command{}

	// Building subcommands for each generator
	for k, v := range gCore.Config.Generators {
		flagSet := []cli.Flag{}
		// Build flagSet for generated help:
		for _, f := range v.Vars { // Vars
			if f.Flag != "" {
				fl := cli.StringFlag{
					Name:  f.Flag,
					Usage: f.Question,
				}
				flagSet = append(flagSet, fl)
			}
		}
		// Build command
		gen := v
		genName := k
		cmd := cli.Command{
			Name:  k,
			Usage: v.Description,
			Flags: flagSet,
			Action: func(cl *cli.Context) error {
				var err error
				log.Println("Source generator:", cl.Command.Name)
				if cl.NArg() == 0 { // One command
					cli.ShowCommandHelp(cl, cl.Command.Name)
					return err
				}
				//var targetFile string
				// Command specific data
				gCore.Data["name"] = cl.Args().Get(0) // Change this to other name

				flagOrAsk(cl, gen.Vars, gCore.Data)
				gCore.Generate(genName)
				return nil
			},
		}
		subCommands = append(subCommands, cmd)
	}
	ret.Subcommands = subCommands

	return ret
}
