package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gohxs/boiler/internal/config"

	"github.com/urfave/cli"
	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"
)

func commandNew() cli.Command {
	return cli.Command{
		Name:      "new",
		ArgsUsage: "[repository/source] [projname]",
		Usage:     "Initialize the boilerplate",
		Action:    actionNew,
	}
}

func actionNew(c *cli.Context) error {

	if c.NArg() < 2 {
		return cli.ShowCommandHelp(c, "init")
	}
	stdin, ok := c.App.Metadata["stdin"].(io.Reader)
	if !ok {
		stdin = os.Stdin
	}
	source := c.Args().Get(0)
	name := c.Args().Get(1)

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

	// Template data

	//////////////////////
	// Method for this
	//////
	cfg, err := config.FromFile(filepath.Join(srcdir, ".boiler", "config.yml"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	udata := map[string]interface{}{} // UserVars
	if cfg != nil {
		// Data questions
		in := bufio.NewReader(stdin)
		for _, v := range cfg.UserVars {
			// User interaction
			fmt.Printf("[%s] %s (%s)? ", v.Name, v.Question, v.Default)
			d, _, _ := in.ReadLine()
			str := string(d)
			if str == "" {
				udata[v.Name] = v.Default
			} else {
				udata[v.Name] = str
			}
		}
	}
	udata["projname"] = name
	// Store data in boiler folder

	// Place this somewhere that loads defaults
	data := map[string]interface{}{} // Merge user data with data
	for k, v := range udata {
		data[k] = v
	}

	// Setup global template vars
	data["projroot"], _ = filepath.Abs(name)
	data["curdir"], _ = os.Getwd() // Not good?

	// Setup vars
	err = ProcessDir(srcdir, name, data)
	if err != nil {
		return err
	}
	fmt.Println("Created project:", name)
	ydata, err := yaml.Marshal(udata)
	if err != nil {
		return err
	}
	//mkdir all .boiler in case it does not exists
	boilerPath := filepath.Join(name, ".boiler")
	os.MkdirAll(boilerPath, os.FileMode(0755)) // ignore error

	err = ioutil.WriteFile(filepath.Join(boilerPath, "user.yml"), ydata, os.FileMode(0644))

	return err
}
