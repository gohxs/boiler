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
)

func commandInit() cli.Command {
	return cli.Command{
		Name:      "init",
		Aliases:   []string{"i"},
		ArgsUsage: "[repository] [projname]",
		Usage:     "Initialize the boilerplate",
		Action:    actionInit,
	}
}
func actionInit(c *cli.Context) error {

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
	data := map[string]interface{}{}

	//////////////////////
	// Method for this
	//////
	cfg, err := config.FromFile(filepath.Join(srcdir, ".boiler", "config.yaml"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if cfg != nil {
		// Data questions
		in := bufio.NewReader(stdin)
		for k, v := range cfg.Vars {
			fmt.Printf("%s? ", v.Question)
			d, _, _ := in.ReadLine()
			data[k] = string(d)
		}
	}

	// Setup global template vars
	data["projroot"], _ = filepath.Abs(name)
	data["curdir"], _ = os.Getwd()
	data["projname"] = name

	// Setup vars
	err = ProcessDir(srcdir, name, data)
	if err != nil {
		return err
	}
	fmt.Println("Created project:", name)

	return nil
}
