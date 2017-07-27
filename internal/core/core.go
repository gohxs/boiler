package core

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/gohxs/boiler/internal/config"
)

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
		return nil, err
	}

	return c, nil

}

func (c *Core) Init() (err error) {
	// Defaults
	c.ConfigFile = filepath.Join(c.ProjRoot, BoilerDir, "config.yml")
	c.UserVarFile = filepath.Join(c.ProjRoot, BoilerDir, "user.yml")

	// Load config
	err = config.FromFile(c.ConfigFile, &c.Config)
	if err != nil { // Ignore error
		//log.Println("Warning: non .boiler project")
		//return err
	}

	// Load vars from user.yml if exists
	userFile, err := ioutil.ReadFile(c.UserVarFile)
	if err == nil {
		yaml.Unmarshal(userFile, c.Data) // Add to data
	}

	c.defaultVars()

	return nil
}

func (c *Core) defaultVars() {
	log.Println("Setting default vars")
	c.Data["curdir"], _ = os.Getwd()
	c.Data["proj-root"] = c.ProjRoot
	// Date tools
	now := time.Now().UTC()
	c.Data["year"] = now.Format("2006")
	c.Data["month"] = now.Format("Jan")
	c.Data["day"] = now.Format("02")
}

func (c *Core) Generate(name string) (err error) {
	gen := c.Config.Generators[name]
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

		log.Println("Ext is :", ext)
		if !strings.HasSuffix(targetFile, ext) {
			targetFile += ext
		}
		srcPath := filepath.Join(c.ProjRoot, ".boiler", "templates", f.Source)
		log.Println("Processing file:", srcPath, " to ", targetFile)

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
