package cmd

import (
	"fmt"

	"github.com/gohxs/boiler/internal/core"
	"github.com/spf13/cobra"
)

// Add Package command
func init() {

	// Build flags here too i guess
	cmd := &cobra.Command{
		Use:   "generators",
		Short: "List generators of the current boilerplate",
		Run: func(cmd *cobra.Command, args []string) {
			for k := range core.Config().Generators {
				fmt.Println(k)
			}
		},
	}
	RootCmd.AddCommand(cmd)
}
