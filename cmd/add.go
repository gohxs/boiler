package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gohxs/boiler/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Add())
}

// Add Package command
func Add() *cobra.Command {

	// Build flags here too i guess
	ret := &cobra.Command{
		Use:   "add [file]",
		Short: "Add a file based on template",
		Long:  "using file as an argument requires a generator containing extension based aliases",
		Run: func(cmd *cobra.Command, args []string) {
			// Not enough args
			if len(args) < 1 {
				cmd.Help()
				return
			}
			// Fetch a file name from this, improve this
			genName := args[0]
			name := ""
			if len(args) != 2 {
				// We load based on extension on arg0
				ext := filepath.Ext(genName)
				if ext != "" {
					name = genName
					genName = ext // Gen command will be ext if available
				}
			} else {
				name = args[1]
			}
			if name == "" {
				cmd.Help()
				//cli.ShowSubcommandHelp(cl)
				return
			}

			gen := core.Cur.GetGenerator(genName)
			if gen == nil {
				cmd.Printf("Generator %s does not exists\n", genName)
				cmd.Help()
				return
			}

			flagOrAsk(cmd, gen.Vars, core.Cur.Data)
			core.Cur.Generate(genName, name)
		},
	}

	// Build sub commands from gen:
	for k, v := range core.Cur.Config.Generators {
		gen := v
		genName := k
		cmd := cobra.Command{
			Use:     fmt.Sprintf("%s [name]", genName),
			Long:    v.Description,
			Aliases: v.Aliases,
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) != 1 {
					cmd.Help()
					return
				}
				flagOrAsk(cmd, gen.Vars, core.Cur.Data)
				core.Cur.Generate(genName, args[0])
			},
		}
		for _, f := range gen.Vars {
			flagDefault, _ := core.ProcessString(f.Default, core.Cur.Data)
			flParts := strings.Split(f.Flag, ",")
			if len(flParts) > 1 {
				cmd.Flags().StringP(strings.TrimSpace(flParts[0]), strings.TrimSpace(flParts[1]), flagDefault, f.Question)
			} else {
				cmd.Flags().String(f.Flag, flagDefault, f.Question)
			}
		}
		//flagFromVars(cmd.Flags(), gen.Vars)
		ret.AddCommand(&cmd)
	}

	return ret
}
