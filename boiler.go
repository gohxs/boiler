//go:generate go get dev.hexasoftware.com/hxs/genversion
//go:generate genversion -package boiler -out version.go
package boiler

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohxs/boiler/config"

	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"
)

const (
	// BoilerExt extension of boiler files
	BoilerExt = ".boiler"
	// BoilerDir directory name inside the boiler project
	BoilerDir = ".boiler"

	VARPROJNAME = "projName"
	VARPROJROOT = "projRoot"
	VARPROJDATE = "projDate"
)

var (
	// Global from currentDir if exists
	global, _ = From(".")

	// Accessors to the global proj (curdir project)

	Generate     = global.Generate
	GetGenerator = global.GetGenerator
	Config       = global.Config
	Data         = global.Data
	ProjRoot     = global.ProjRoot
	Name         = global.Name
)

// Core This is a projConfig/Proj
type Core struct {
	name        string
	config      config.Config
	data        map[string]interface{}
	projRoot    string
	configFile  string
	userVarFile string
	isTemporary bool // If temporary such as git temporary sources, remove
}

// New instantiate a dumb core
func New(path string) *Core {
	return &Core{projRoot: path, data: map[string]interface{}{}}
}

//Init initializes core based on config
func (c *Core) Init() (err error) {
	// Defaults
	c.configFile = filepath.Join(c.projRoot, BoilerDir, "config.yml")
	c.userVarFile = filepath.Join(c.projRoot, BoilerDir, "user.yml")

	// Load config
	err = config.FromFile(c.configFile, &c.config)
	if err != nil && !os.IsNotExist(err) { // Ignore error if does not exists
		return err
		//return err
	}
	// Load vars from user.yml if exists
	userFile, err := ioutil.ReadFile(c.userVarFile)
	if err == nil { // NO ERROR intentional, we only unmarshal if file exists else its ok to go on
		yaml.Unmarshal(userFile, c.data) // Add to data
	}

	if projName, ok := c.data[VARPROJNAME]; ok {
		c.name = projName.(string)
	}
	return nil
}

// CloneTo and process the boiler plate to destination
func (c *Core) CloneTo(dest string) (err error) {
	name := filepath.Base(dest)
	dir := filepath.Dir(dest)
	if dir != "" {
		err := os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	c.name = name

	c.data[VARPROJNAME] = name
	c.data[VARPROJDATE] = time.Now().UTC()
	//fmt.Print("Generating project...\n\n")

	// Setup vars
	err = ProcessPath(c.projRoot, dest, c.data)
	if err != nil {
		return err
	}
	ydata, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}
	//mkdir all .boiler in case it does not exists
	boilerPath := filepath.Join(dest, ".boiler")
	os.MkdirAll(boilerPath, os.FileMode(0755)) // ignore error

	err = ioutil.WriteFile(filepath.Join(boilerPath, "user.yml"), ydata, os.FileMode(0644))

	return nil

}

//////////////////////////////////
// GENERATORS
////////////////////////

// GetGenerator Fetches generator by name/alias
func (c *Core) GetGenerator(name string) *config.Generator {

	for k, v := range c.config.Generators {
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

// Generate uses a generator
func (c *Core) Generate(generator string, name string) (err error) {

	// DefaultVars here?
	c.data["name"] = name // Name or target
	c.data[VARPROJROOT] = c.ProjRoot()
	c.data["time"] = time.Now().UTC() //curTime
	c.data["curdir"], _ = os.Getwd()  //currentDir (useful for file paths in config)?

	gen := c.GetGenerator(generator)
	// Each file
	for _, f := range gen.Files {
		targetFile, err := ProcessString(f.Target, c.data)
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
		srcPath := filepath.Join(c.projRoot, ".boiler", "templates", f.Source)
		fmt.Println("Generating file:", targetFile)

		// Create dir
		dir, _ := filepath.Split(targetFile)
		os.MkdirAll(dir, os.FileMode(0755))
		err = ProcessPath(srcPath, targetFile, c.data)
		if err != nil {
			return err
		}
	}
	return nil

}
func (c *Core) Close() {
	if c.isTemporary {
		defer os.RemoveAll(c.projRoot)

	}
}

//////////////////////
// GETTERS
////////

// Config returns configuration
func (c *Core) Config() *config.Config {
	return &c.config
}

// Data returns Data from core
func (c *Core) Data() map[string]interface{} {
	return c.data
}

// ProjRoot returns project path in filesystem
func (c *Core) ProjRoot() string {
	return c.projRoot
}

// Name returns proj name
func (c *Core) Name() string {
	return c.name
}

////////////////////////////////
// Specialized factory
////////////////

// From Should be multi purpose
func From(source string) (*Core, error) {
	var (
		srcdir string
		err    error
	)
	if source == "" { // Special case if empty load from current dir?
		srcdir = "." //, _ = os.Getwd()
	} else {
		srcdir = source
		u, err := url.Parse(source)
		if err != nil {
			return nil, err
		}
		// Git to tmpdir // Maybe move this to cmd
		if u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "git" {
			srcdir, err = ioutil.TempDir(os.TempDir(), "boiler")
			if err != nil {
				return nil, err
			}
			_, err = git.PlainClone(srcdir, false, &git.CloneOptions{
				URL:      source,
				Progress: os.Stdout,
			})
			if err != nil {
				return nil, err
			}
			c := New(srcdir)
			c.isTemporary = true
			err = c.Init()
			if err != nil {
				c.Close() // Close if init error
				return nil, err
			}
			return c, nil
		}
	}
	// Solve dir into .boiler root
	srcdir, err = solveProjRoot(srcdir)
	if err != nil {
		return nil, err
	}
	// check if exists
	_, err = os.Stat(srcdir)
	if err != nil {
		return nil, err
	}
	// TempCore
	c := New(srcdir)
	err = c.Init()
	return c, err
}
