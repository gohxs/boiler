package cmd

import (
	"fmt"

	"github.com/gohxs/boiler"
	"github.com/spf13/cobra"
)

// Add Package command
func init() {

	// Build flags here too i guess
	cmd := &cobra.Command{
		Use:   "generators",
		Short: "List generators of the current boilerplate",
		Run: func(cmd *cobra.Command, args []string) {
			for k := range boiler.Config().Generators {
				fmt.Println(k)
			}
		},
	}
	RootCmd.AddCommand(cmd)
}
