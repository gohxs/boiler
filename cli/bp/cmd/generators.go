package cmd

import (
	"fmt"
	"os"

	"github.com/gohxs/boiler"
	"github.com/spf13/cobra"
)

// Add Package command
func init() {

	b := Boiler()
	cmd := &cobra.Command{
		Use:   "generators",
		Short: "generators related commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List generators of the current project",
		Run: func(cmd *cobra.Command, args []string) {
			for k := range b.Config.Generators {
				fmt.Println(k) // stdout
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "fetch [boilerplate] [generator name] <local name>",
		Short: "Fetch generators from other boilerplates",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) < 2 {
				cmd.Help()
				return
			}
			if !b.IsBoiler {
				cmd.Printf("%s is not a boiler project\n", b.ProjRoot)
				cmd.Printf("Run '%s init [name]' to initialize boiler in current dir\n", os.Args[0])
				return
			}
			cmd.Println("Boiler root:", b.ProjRoot)

			repo := args[0]
			genName := args[1]
			localName := genName // Default to genName
			if len(args) == 3 {
				localName = args[2]
			}

			// possible move this logic to boiler instead of command

			if gen := b.GetGenerator(localName); gen != nil {
				cmd.Println("Generator already exists")
				return
			}
			secProj, err := boiler.From(repo)
			if err != nil {
				cmd.Println(err)
				return
			}
			defer secProj.Close()
			err = b.GeneratorFetch(secProj, genName, localName)
			if err != nil {
				cmd.Println(err)
				return
			}
			cmd.Printf("Generator '%s' added\n", genName)

		},
	})

	RootCmd.AddCommand(cmd)
}
