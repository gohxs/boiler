package cliapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// TODO: Reduce this to single command
func commandCreate() cli.Command {
	ret := cli.Command{
		Name:    "create",
		Aliases: []string{"c", "Bigfstring"},
		Usage:   "Create a thing from template",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "vars",
				Usage: "Show the available flags for [generator]",
			},
		},

		Action: func(cl *cli.Context) error {
			if cl.NArg() < 1 {
				cli.ShowSubcommandHelp(cl)
				return nil
			}
			genName := cl.Args().Get(0)
			name := ""
			ext := ""
			if cl.NArg() != 2 {
				// We load based on extension on arg0
				ext = filepath.Ext(genName)
				if ext != "" {
					name = genName
					genName = ext // Gen command will be ext if available
				}
			} else {
				name = cl.Args().Get(1)
			}

			//fmt.Printf("Generator: %s, name: %s\n", genName, name)
			// Based on command we fetch the right generator including wildcards
			gen := gCore.GetGenerator(genName)
			if gen == nil {
				cli.ShowSubcommandHelp(cl)
				return fmt.Errorf("Generator %s does not exists", genName)
			}

			if cl.Bool("vars") {
				fmt.Printf("VARS: %s\n", genName)
				// Show flags here
				flagSet := flagFromVars(gen.Vars)
				for _, fl := range flagSet {
					fmt.Printf("  %s\n", fl)
				}
				return nil
			}
			if name == "" {
				cli.ShowSubcommandHelp(cl)
				return fmt.Errorf("Missing [name] parameter")
			}

			flagOrAsk(cl, gen.Vars, gCore.Data)
			gCore.Generate(genName, name)

			// Arg1 can be a generator or extension based generator
			return nil
		},
		UsageText: "Several usage entries including flag help",
	}

	ret.UsageText = ret.UsageText + "\n\nGENERATORS:"

	for k, v := range gCore.Config.Generators {
		genName := k
		if len(v.Aliases) > 0 {
			genName = fmt.Sprintf("%s, %s", k, strings.Join(v.Aliases, ", "))
		}
		// Build command
		ret.UsageText += fmt.Sprintf("\n     %s\t%s", genName, v.Description)

		/*flagSet := flagFromVars(v.Vars)
		for _, fl := range flagSet {
			ret.UsageText += fmt.Sprintf("        %s\n", fl)
		}*/

	}

	// Build UsageText, show available commands
	//cli.DefaultAppComplete(

	//ret.Subcommands = subCommands

	return ret
}
