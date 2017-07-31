package cmd

import (
	"github.com/spf13/cobra"
)

// Add Package command
func init() {
	b := Boiler()
	// Build flags here too i guess
	cmd := &cobra.Command{
		Use:     "init [name]",
		Aliases: []string{"i"},
		Short:   "Initialize .boiler in current dir ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			name := args[0]

			err := b.InitProj(name)
			if err != nil {
				cmd.Println(err)
			}

			cmd.Println("Created:", b.ConfigFile)

			cmd.Printf("Project '%s' Created\n", name)

			// Find if we have a proj in parent dir
			// Move this to boiler
		},
	}
	RootCmd.AddCommand(cmd)
}
