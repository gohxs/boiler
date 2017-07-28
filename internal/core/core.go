package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"

	"github.com/gohxs/boiler/internal/config"
)

const (
	// BoilerExt extension of boiler files
	BoilerExt = ".boiler"
	// BoilerDir directory name inside the boiler project
	BoilerDir = ".boiler"
)

var (
	// Global from currentDir if exists
	global, _ = FromCurDir()

	// Accessors to the global proj (curdir project)

	Generate     = global.Generate
	GetGenerator = global.GetGenerator
	Config       = global.Config
	Data         = global.Data
	ProjRoot     = global.ProjRoot
	Name         = global.Name
)

func init() {
	//global.DefaultVars() // Make it create default vars
}

// This is a projConfig/Proj
type Core struct {
	name        string
	config      config.Config
	data        map[string]interface{}
	projRoot    string
	configFile  string
	userVarFile string
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
	return nil
}

// DefaultVars Set vars
/*func (c *Core) DefaultVars() {
	c.data["curdir"], _ = os.Getwd()  //?
	c.data["time"] = time.Now().UTC() // curTime

	c.data["projRoot"] = c.projRoot //?
	// Date tools
}*/

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

	c.data["projName"] = name
	c.data["projDate"] = time.Now().UTC()
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
	c.data["name"] = name

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
	srcdir := source

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
	// TempCore
	c := New(srcdir)

	err = c.Init()
	return c, err
}

// FromCurDir returns a boiler from current path
func FromCurDir() (*Core, error) {
	cwd, _ := os.Getwd()
	return SearchPath(cwd)
}

// SearchPath Find project in parent folders
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
