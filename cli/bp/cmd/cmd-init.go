package cmd

import (
	"github.com/gohxs/boiler"
	"github.com/spf13/cobra"
)

// Add Package command
func init() {
	//b := Boiler()
	// Build flags here too i guess
	cmd := &cobra.Command{
		Use:     "init [name] <path>",
		Aliases: []string{"i"},
		Short:   "Initialize .boiler in current dir ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			name := args[0]
			bpath := "."
			if len(args) > 1 {
				bpath = args[1]
			}

			bProj := boiler.New(bpath)
			if err := bProj.Init(); err != nil {
				cmd.Println(err)
				return
			}
			if err := bProj.InitProj(name); err != nil {
				cmd.Println(err)
				return
			}

			cmd.Println("Created:", bProj.ConfigFile)

			cmd.Printf("Project '%s' Created\n", name)

			// Find if we have a proj in parent dir
			// Move this to boiler
		},
	}
	RootCmd.AddCommand(cmd)
}
