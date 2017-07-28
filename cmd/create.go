package cmd

import (
	"github.com/gohxs/boiler/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Create())
}

// Create command creates clones a template into destination directory
func Create() *cobra.Command {
	// Transform flags
	ret := &cobra.Command{
		Use:  "create [repository/source] [projname]",
		Long: "Create new project from a boilerplate",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			source := args[0]
			dest := args[1]

			_, err := core.NewProj(source, dest, func(c *core.Core) {
				flagOrAsk(cmd, c.Config.UserVars, c.Data)
			})
			if err != nil {
				cmd.Println("ERR:", err)
			}

		},
	}
	return ret
}
