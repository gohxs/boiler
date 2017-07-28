package cliapp

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gohxs/boiler/internal/core"
	"github.com/urfave/cli"
	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"
)

func commandNew() cli.Command {
	// Transform flags
	ret := cli.Command{
		Name:      "new",
		ArgsUsage: "[repository/source] [projname]",
		Usage:     "Create new Project from a boilerplate",
		Action:    actionNew,
	}
	return ret
}

func actionNew(cl *cli.Context) error {

	if cl.NArg() < 2 {
		return cli.ShowCommandHelp(cl, cl.Command.Name)
	}
	source := cl.Args().Get(0)
	name := cl.Args().Get(1)

	srcdir := source
	u, err := url.Parse(source)
	if err != nil {
		return err
	}
	// Git to tmpdir
	if u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "git" {
		srcdir, err = ioutil.TempDir(os.TempDir(), "boiler")
		if err != nil {
			return err
		}
		defer os.RemoveAll(srcdir)
		_, err = git.PlainClone(srcdir, false, &git.CloneOptions{
			URL:      source,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(srcdir) // Check if source exists
	if err != nil {
		return err

	}

	c := core.New(srcdir)
	// Template data
	err = c.Init()
	if err != nil {
		return err
	}
	fmt.Println(c.Config.Description)
	fmt.Println("-----")

	udata := map[string]interface{}{} // UserVars
	// User defined param
	// Creation time params
	udata["projName"] = name
	udata["projDate"] = time.Now().UTC()
	// Store data in boiler folder
	// Attempt
	flagOrAsk(cl, c.Config.UserVars, udata)

	// Merge map
	for k, v := range udata {
		c.Data[k] = v
	}

	fmt.Print("Generating project...\n\n")
	// Setup vars
	err = core.ProcessPath(srcdir, name, c.Data)
	if err != nil {
		return err
	}
	fmt.Println("Created project:", name)
	ydata, err := yaml.Marshal(c.Data)
	if err != nil {
		return err
	}
	//mkdir all .boiler in case it does not exists
	boilerPath := filepath.Join(name, ".boiler")
	os.MkdirAll(boilerPath, os.FileMode(0755)) // ignore error

	err = ioutil.WriteFile(filepath.Join(boilerPath, "user.yml"), ydata, os.FileMode(0644))

	return err
}
