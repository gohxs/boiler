package cliapp

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/gohxs/boiler/internal/config"
	"github.com/urfave/cli"
)

// TODO: Reduce this to single command
func commandCreatebak() cli.Command {

	ret := cli.Command{
		Name:  "createbak",
		Usage: "Create a thing from template",
		Before: func(cl *cli.Context) error {
			log.Println("Select which command to run")
			log.Println("Does this works?")
			log.Println("Args:", cl.Args())

			if cl.NArg() > 1 {
				return nil
			}

			genName := cl.Args().Get(0)
			ext := filepath.Ext(genName)
			// We load based on extension on arg0
			if ext == "" {
				return nil
			}
			name := genName
			genName = ext // Gen command will be ext if available
			newArgs := []string{os.Args[0], cl.Parent().Command.Name, genName, name, "-h"}
			log.Println("Rerun app with proper extensions", newArgs)
			// We had to change things
			cl.App.Run(newArgs) // Rerun App with New changed arguments

			return errors.New("") //Sort circuit

		},
	}

	subCommands := []cli.Command{}
	// Building subcommands for each generator
	for k, v := range gCore.Config.Generators {
		flagSet := flagFromVars(v.Vars)
		// Build command
		gen := v
		genName := k
		cmd := cli.Command{
			Name:    k,
			Aliases: v.Aliases,
			Usage:   v.Description,
			Flags:   flagSet,
			Action: func(cl *cli.Context) error {
				var err error
				if cl.NArg() != 1 { // One command
					cli.ShowCommandHelp(cl, cl.Command.Name)
					return err
				}
				//var targetFile string
				// Command specific data

				flagOrAsk(cl, gen.Vars, gCore.Data)
				gCore.Generate(genName, cl.Args().Get(0))
				return nil
			},
		}
		subCommands = append(subCommands, cmd)
	}
	ret.Subcommands = subCommands

	//ret.Subcommands = subCommands

	return ret
}

func flagFromVars(vars []config.UserVar) []cli.Flag {
	flagSet := []cli.Flag{}
	// Build flagSet for generated help:
	for _, f := range vars { // Vars
		if f.Flag != "" {
			fl := cli.StringFlag{
				Name:  f.Flag,
				Usage: f.Question,
			}
			flagSet = append(flagSet, fl)
		}
	}
	return flagSet
}
