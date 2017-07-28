package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"

	"github.com/gohxs/boiler/internal/config"
)

var (
	Cur, _ = FromCurDir()
)

func init() {
	Cur.DefaultVars() // Make it create default vars
}

// This is a projConfig/Proj
type Core struct {
	Config      config.Config
	Data        map[string]interface{}
	ProjRoot    string
	ConfigFile  string
	UserVarFile string
}

const (
	BoilerExt = ".boiler"
	BoilerDir = ".boiler"
)

// Create a new proj
func NewProj(source, dest string, onInit func(c *Core)) (*Core, error) {
	srcdir := source
	name := dest

	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	// Git to tmpdir
	if u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "git" {
		srcdir, err = ioutil.TempDir(os.TempDir(), "boiler")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(srcdir)
		_, err = git.PlainClone(srcdir, false, &git.CloneOptions{
			URL:      source,
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, err
		}
	}
	_, err = os.Stat(srcdir) // Check if source exists
	if err != nil {
		return nil, err
	}

	c := New(srcdir)
	// Template data
	err = c.Init()
	if err != nil {
		return nil, err
	}
	fmt.Println(c.Config.Description)
	fmt.Println("-----")

	onInit(c) // Execute what we want on src

	/*udata := map[string]interface{}{} // UserVars
	// User defined param
	// Creation time params
	udata["projName"] = name
	udata["projDate"] = time.Now().UTC()
	// Store data in boiler folder
	// Attempt

	// Merge map
	for k, v := range udata {
		c.Data[k] = v
	}*/
	c.Data["projName"] = name
	c.Data["projData"] = time.Now().UTC()

	fmt.Print("Generating project...\n\n")
	// Setup vars
	err = ProcessPath(srcdir, name, c.Data)
	if err != nil {
		return nil, err
	}
	fmt.Println("Created project:", name)
	ydata, err := yaml.Marshal(c.Data)
	if err != nil {
		return nil, err
	}
	//mkdir all .boiler in case it does not exists
	boilerPath := filepath.Join(name, ".boiler")
	os.MkdirAll(boilerPath, os.FileMode(0755)) // ignore error

	err = ioutil.WriteFile(filepath.Join(boilerPath, "user.yml"), ydata, os.FileMode(0644))

	return c, nil

}

func New(path string) *Core {
	c := &Core{Data: map[string]interface{}{}}
	c.ProjRoot = path
	// LoadVars
	return c
}

// Factories
func FromCurDir() (*Core, error) {
	cwd, _ := os.Getwd()
	return SearchPath(cwd)
}

func SearchPath(path string) (*Core, error) {
	projRoot := solveProjRoot(path)
	if projRoot == "" {
		projRoot = path
	}
	c := New(projRoot)
	err := c.Init()
	if err != nil {
		log.Println("Err:", err)
		return c, err
	}

	return c, nil

}

func (c *Core) Init() (err error) {
	// Defaults
	c.ConfigFile = filepath.Join(c.ProjRoot, BoilerDir, "config.yml")
	c.UserVarFile = filepath.Join(c.ProjRoot, BoilerDir, "user.yml")

	// Load config
	err = config.FromFile(c.ConfigFile, &c.Config)
	if err != nil && !os.IsNotExist(err) { // Ignore error if does not exists
		return err
		//return err
	}

	// Load vars from user.yml if exists
	userFile, err := ioutil.ReadFile(c.UserVarFile)
	if err == nil { // NO ERROR intentional, we only unmarshal if file exists else its ok to go on
		yaml.Unmarshal(userFile, c.Data) // Add to data
	}
	return nil
}

func (c *Core) DefaultVars() {
	c.Data["curdir"], _ = os.Getwd()
	c.Data["projRoot"] = c.ProjRoot
	// Date tools
	c.Data["time"] = time.Now().UTC() // curTime
}

func (c *Core) GetGenerator(name string) *config.Generator {

	for k, v := range c.Config.Generators {
		if k == name {
			return &v // Is a copy of?
		}
		for _, al := range v.Aliases {
			if al == name {
				return &v
			}
		}
	}
	return nil // not found
}

func (c *Core) Generate(generator string, name string) (err error) {

	c.Data["name"] = name

	gen := c.GetGenerator(generator)
	// Each file
	for _, f := range gen.Files {
		targetFile, err := ProcessString(f.Target, c.Data)
		if err != nil {
			return err
		}
		ext := filepath.Ext(f.Source)
		if ext == ".boiler" {
			ext = filepath.Ext(f.Source[:len(f.Source)-7])
		}

		//log.Println("Ext is :", ext)
		if !strings.HasSuffix(targetFile, ext) {
			targetFile += ext
		}
		srcPath := filepath.Join(c.ProjRoot, ".boiler", "templates", f.Source)
		fmt.Println("Generating file:", targetFile)

		// Create dir
		dir, _ := filepath.Split(targetFile)
		os.MkdirAll(dir, os.FileMode(0755))
		err = ProcessPath(srcPath, targetFile, c.Data)
		if err != nil {
			return err
		}
	}
	return nil

}

// Find project root from start path
func solveProjRoot(start string) string {
	cwd := filepath.Clean(start)
	for f := cwd; f != ""; cwd, f = path.Split(cwd) {
		cwd = filepath.Clean(cwd)
		//log.Println("Cwd", cwd)
		boilerpath := filepath.Join(cwd, ".boiler")
		st, err := os.Stat(boilerpath)
		if err == nil && st.IsDir() { // ignore error
			return cwd
		}
	}
	return ""
}
