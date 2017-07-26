package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/gohxs/boiler/internal/config"
	"github.com/urfave/cli"
)

func commandCreate() cli.Command {

	subCommands := []cli.Command{}
	cwd, _ := os.Getwd()
	projroot := solveProjRoot(cwd)
	if projroot != "" {
		cfg, err := config.FromFile(filepath.Join(projroot, ".boiler", "config.yml"))
		if err != nil {
			panic(err)
		}
		// Create a data loader else where
		data := map[string]interface{}{}
		userFile, err := ioutil.ReadFile(filepath.Join(projroot, ".boiler", "user.yml"))
		if err == nil {
			yaml.Unmarshal(userFile, data)
		}
		data["curdir"] = cwd
		data["projroot"] = projroot

		// Building subcommands
		log.Println("Sub things:", data)
		log.Println("Loading generators:", cfg.Generators)
		for k, v := range cfg.Generators {
			flagSet := []cli.Flag{}
			// Build flag:
			for _, f := range v.Flags {
				log.Println("Adding flag:", f)
				fl := cli.StringFlag{
					Name:  f,
					Usage: fmt.Sprintf("Flag: %s sets %s param", f, f),
				}
				flagSet = append(flagSet, fl)
			}
			cmd := cli.Command{
				Name:  k,
				Usage: v.Target,
				Flags: flagSet,
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						cli.ShowCommandHelp(c, k)
						return err
					}
					var targetFile string
					data["name"] = c.Args().Get(0)
					for _, f := range c.FlagNames() {
						if !c.IsSet(f) {
							continue // Is required? or ask?
						}
						log.Println("Loding value for:", f)
						data[f] = c.String(f)
					}
					{
						buf := bytes.NewBuffer([]byte{})
						t, err := template.New("l").Parse(v.Target)
						if err != nil {
							return err
						}
						t.Execute(buf, data)
						targetFile = buf.String()
					}
					if !strings.HasSuffix(targetFile, v.Ext) {
						targetFile += "." + v.Ext
					}
					log.Println("Create file:", targetFile)
					// Read boiler file
					// check if dir or file
					ProcessFile(filepath.Join(projroot, ".boiler", "templates", v.Source), targetFile, data)

					return nil
				},
			}
			subCommands = append(subCommands, cmd)
		}
	}
	log.Printf("Sub commands: %v", subCommands)

	return cli.Command{
		Name:        "create",
		Usage:       "Create a thing from template",
		Subcommands: subCommands,
	}
}

func actionCreate(c *cli.Context) error {
	// Find boiler root
	return nil
}

func solveProjRoot(start string) string {
	cwd := path.Clean(start)
	for f := cwd; f != ""; cwd, f = path.Split(cwd) {
		cwd = path.Clean(cwd)
		log.Println("Cwd", cwd)
		boilerpath := path.Join(cwd, ".boiler")
		st, err := os.Stat(boilerpath)
		if err == nil && st.IsDir() { // ignore error
			return cwd
		}

	}
	return ""
}
