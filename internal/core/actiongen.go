package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/urfave/cli"
)

func commandGenerate() cli.Command {
	return cli.Command{
		Name:    "generate",
		Aliases: []string{"new", "g"},
		Usage:   "Generate a thing",
		Action:  actionGenerate,
	}
}
func actionGenerate(c *cli.Context) error {
	// Find boiler root
	cwd, _ := os.Getwd()
	projroot := solveProjRoot(cwd)
	if projroot == "" {
		return errors.New("Not a boiler project")
	}
	log.Println("Project root:", projroot)

	// Do something

	fmt.Println("Generate will Grab the arguments")
	fmt.Println("Look for config yaml in ../[../../]")
	fmt.Println("Show all the options in config.yaml")
	return nil
}

func solveProjRoot(start string) string {
	cwd := start

	for ; cwd != ""; cwd, _ = path.Split(cwd) {
		boilerpath := path.Join(cwd, ".boiler")
		st, err := os.Stat(boilerpath)
		if err == nil && st.IsDir() { // ignore error
			return cwd
		}

	}
	return ""
}
