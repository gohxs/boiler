package cmd

import (
	"os"

	"github.com/gohxs/boiler"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:     "new [repository/source] [projname]",
		Aliases: []string{"n", "create"},
		Short:   "Create new project from a boilerplate",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			source := args[0]
			dest := args[1]
			// Check existent destination
			_, err := os.Stat(dest)
			if !os.IsNotExist(err) {
				cmd.Printf("%s already exists\n", dest)
				return
			}

			cmd.Printf("Loading boilerplate from %s\n", source)
			c, err := boiler.From(source)
			if err != nil {
				cmd.Println(err)
				return
			}
			defer c.Close()

			cmd.Println(c.Config.Description)
			cmd.Println("-----")
			// Set vars on cur Plate
			flagOrAsk(cmd, c.Config.UserVars, c.Data)

			cmd.Print("Generating project...\n")
			err = c.CloneTo(dest)
			if err != nil {
				cmd.Println(err)
				return
			}

			cmd.Println("Created project:", dest, c.Data["projName"])

		},
	}
	RootCmd.AddCommand(cmd)
}
