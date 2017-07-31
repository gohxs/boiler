package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gohxs/boiler"
	"github.com/gohxs/boiler/config"
	"github.com/spf13/cobra"
)

// Add Package command
func init() {

	b := Boiler()
	// Build flags here too i guess
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
			cmd.Println("Boiler project:", b.ProjRoot)

			repo := args[0]
			genName := args[1]
			localName := genName // Default to genName
			if len(args) == 3 {
				localName = args[2]
			}

			if _, ok := b.Config.Generators[localName]; ok {
				cmd.Println("Generator already exists")
				return
			}

			secProj, err := boiler.From(repo)
			if err != nil {
				panic(err)
			}
			defer secProj.Close()

			gen := secProj.GetGenerator(genName)
			if gen == nil {
				cmd.Println("Generator does not exists in requested boilerplate")
				return
			}
			// Should we copy global vars to local project aswell???
			// Go trough global vars from secProj and ask to user

			// Reask some vars
			flagOrAsk(cmd, secProj.Config.UserVars, b.Data)

			newGen := config.Generator{}
			newGen.Aliases = gen.Aliases
			newGen.Description = gen.Description
			newGen.Vars = gen.Vars

			// Create local generator entry

			dirPrefix := time.Now().UTC().Format("150402012006")
			dstPath := filepath.Join(b.ProjRoot, boiler.BoilerDir, "templates", dirPrefix)

			for _, f := range gen.Files {
				fsrc := filepath.Join(secProj.ProjRoot, boiler.BoilerDir, "templates", f.Source)
				fdst := filepath.Join(dstPath, f.Source)

				dstDir := filepath.Dir(fdst)
				err = os.MkdirAll(dstDir, os.FileMode(0755)) // ignore error?
				if err != nil {
					cmd.Println(err)
					return
				}

				err = boiler.CopyFile(fsrc, fdst)
				if err != nil {
					cmd.Println(err)
					return
				}
				// Also check if exists
				newGen.Files = append(newGen.Files, config.FileTarget{Source: filepath.Join(dirPrefix, f.Source), Target: f.Target})
			}
			b.Config.Generators[localName] = newGen // entry created

			err = b.Save()
			cmd.Printf("Generator '%s' added\n", genName)

			cmd.Println("\n\n")

		},
	})

	RootCmd.AddCommand(cmd)
}
