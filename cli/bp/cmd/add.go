package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gohxs/boiler"
	"github.com/spf13/cobra"
)

// Add Package command
func init() {

	// Build flags here too i guess
	cmd := &cobra.Command{
		Use:     "add [file]",
		Aliases: []string{"a"},
		Short:   "Add a file based on boilerplate generator",
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

			gen := boiler.GetGenerator(genName)
			if gen == nil {
				cmd.Printf("Generator %s does not exists\n\n", genName)
				cmd.Help()
				return
			}

			flagOrAsk(cmd, gen.Vars, boiler.Data())
			err := boiler.Generate(genName, name)
			if err != nil {
				cmd.Println(err)
			}
		},
	}
	RootCmd.AddCommand(cmd)

	// Flag for listing

	// Build sub commands from gen:
	for k, v := range boiler.Config().Generators {
		gen := v
		genName := k
		subCmd := cobra.Command{
			Use:     fmt.Sprintf("%s [name]", genName),
			Long:    v.Description,
			Aliases: v.Aliases,
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) != 1 {
					cmd.Help()
					return
				}

				flagOrAsk(cmd, gen.Vars, boiler.Data())
				err := boiler.Generate(genName, args[0])
				if err != nil {
					cmd.Println(err)
				}
			},
		}
		for _, f := range gen.Vars {
			flagDefault, _ := boiler.ProcessString(f.Default, boiler.Data())
			flParts := strings.Split(f.Flag, ",")
			if len(flParts) > 1 {
				subCmd.Flags().StringP(strings.TrimSpace(flParts[0]), strings.TrimSpace(flParts[1]), flagDefault, f.Question)
			} else {
				subCmd.Flags().String(f.Flag, flagDefault, f.Question)
			}
		}
		//flagFromVars(cmd.Flags(), gen.Vars)
		cmd.AddCommand(&subCmd)
	}

}
