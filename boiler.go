package boiler

//go:generate go get dev.hexasoftware.com/hxs/genversion
//go:generate genversion -package boiler -out version.go

import (
	"errors"
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
	//VARPROJNAME project name data key
	VARPROJNAME = "projName"
	//VARPROJROOT project root data key
	VARPROJROOT = "projRoot"
	//VARPROJDATE project init date data key
	VARPROJDATE = "projDate"
)

// Core This is a projConfig/Proj
type Core struct {
	Name        string
	Config      config.Config
	Data        map[string]interface{}
	ProjRoot    string
	ConfigFile  string
	UserVarFile string

	IsBoiler    bool
	IsTemporary bool // If temporary such as git temporary sources remove
}

// New instantiate a dumb core
func New(path string) *Core {
	return &Core{ProjRoot: path, Data: map[string]interface{}{}}
}

//Init initializes core based on config
func (c *Core) Init() (err error) {
	// Defaults
	c.ConfigFile = filepath.Join(c.ProjRoot, BoilerDir, "config.yml")
	c.UserVarFile = filepath.Join(c.ProjRoot, BoilerDir, "user.yml")
	c.IsBoiler = true

	// Load config
	err = config.FromFile(c.ConfigFile, &c.Config)
	if os.IsNotExist(err) {
		c.IsBoiler = false
	} else if err != nil { // Ignore error if does not exists
		return err
		//return err
	}
	// Load vars from user.yml if exists
	userFile, err := ioutil.ReadFile(c.UserVarFile)
	if err == nil { // NO ERROR intentional, we only unmarshal if file exists else its ok to go on
		err = yaml.Unmarshal(userFile, c.Data) // Add to data
		if err != nil {
			return err
		}
	}

	if projName, ok := c.Data[VARPROJNAME]; ok {
		c.Name = projName.(string)
	}
	return nil
}

// InitProj create a .boiler and .boiler/config.yml in current path?
func (c *Core) InitProj(name string) (err error) {
	if c.IsBoiler {
		return errors.New("Project already exists")
	}
	boilerPath := filepath.Join(c.ProjRoot, ".boiler")
	err = os.Mkdir(boilerPath, os.FileMode(0755))
	if err != nil {
		return err
	}

	c.Data[VARPROJNAME] = name
	c.Data[VARPROJDATE] = time.Now().UTC()

	return c.Save()
}

// CloneTo and process the boiler plate to destination
func (c *Core) CloneTo(dest string) (err error) {
	name := filepath.Base(dest)
	dir := filepath.Dir(dest)
	if dir != "" {
		err = os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	c.Name = name

	c.Data[VARPROJNAME] = name
	c.Data[VARPROJDATE] = time.Now().UTC()
	//fmt.Print("Generating project...\n\n")

	// TEMPLATE it will copy files
	err = ProcessPath(c.ProjRoot, dest, c.Data)
	if err != nil {
		return err
	}
	newBoiler, err := From(dest)
	if err != nil {
		return err
	}
	newBoiler.Data = c.Data // Clone new Data

	return newBoiler.Save()

}

// Save config and user data
func (c *Core) Save() (err error) {
	// Save config
	err = config.SaveFile(c.ConfigFile, &c.Config)
	if err != nil {
		return err
	}

	ydata, err := yaml.Marshal(c.Data)
	if err != nil {
		return err
	}
	boilerPath := filepath.Join(c.ProjRoot, ".boiler")
	err = os.MkdirAll(boilerPath, os.FileMode(0755)) // ignore error?
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(boilerPath, "user.yml"), ydata, os.FileMode(0644))
	return err

}

//SaveData saves .Data vars
func (c *Core) SaveData() (err error) {
	ydata, err := yaml.Marshal(c.Data)
	if err != nil {
		return err
	}
	boilerPath := filepath.Join(c.ProjRoot, ".boiler")
	err = os.MkdirAll(boilerPath, os.FileMode(0755)) // ignore error
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(boilerPath, "user.yml"), ydata, os.FileMode(0644))
	return err
}

//////////////////////////////////
// GENERATORS
////////////////////////

// GetGenerator Fetches generator by name/alias
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

// Generate uses a generator and clone/process any source files
func (c *Core) Generate(generator string, name string) (err error) {

	// DefaultVars here?
	c.Data["name"] = name // Name or target
	c.Data[VARPROJROOT] = c.ProjRoot
	c.Data["time"] = time.Now().UTC() //curTime
	c.Data["curdir"], _ = os.Getwd()  //currentDir (useful for file paths in config)?

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
		err = os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			return err
		}

		err = ProcessPath(srcPath, targetFile, c.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

//GeneratorFetch fetches a generator from other repository
func (c *Core) GeneratorFetch(srcBoiler *Core, genName, localName string) (err error) {
	if _, ok := c.Config.Generators[localName]; ok {
		return fmt.Errorf("Generator '%s' already exists", localName)
	}

	/*srcProj, err := From(srcBoiler)
	if err != nil {
		return err
	}*/

	if !srcBoiler.IsBoiler {
		return fmt.Errorf("Source '%s' is not a boiler project", srcBoiler.ProjRoot)
	}
	gen := srcBoiler.GetGenerator(genName)
	if gen == nil {
		return fmt.Errorf("Generator '%s' does not exists in '%s'", genName, srcBoiler.ProjRoot)
	}

	newGen := config.Generator{}
	newGen.Aliases = gen.Aliases // This might conflict with existent
	newGen.Description = gen.Description
	newGen.Vars = gen.Vars

	// Create local generator entry

	dirPrefix := time.Now().UTC().Format("20060102150405")
	dstPath := filepath.Join(c.ProjRoot, BoilerDir, "templates", dirPrefix)

	for _, f := range gen.Files {
		fsrc := filepath.Join(srcBoiler.ProjRoot, BoilerDir, "templates", f.Source)
		fdst := filepath.Join(dstPath, f.Source)

		dstDir := filepath.Dir(fdst)
		err = os.MkdirAll(dstDir, os.FileMode(0755)) // ignore error?
		if err != nil {
			break
		}

		err = CopyFile(fsrc, fdst)
		if err != nil {
			break
		}
		// Add entry with generated prefix
		newGen.Files = append(newGen.Files, config.FileTarget{Source: filepath.Join(dirPrefix, f.Source), Target: f.Target})
	}
	if err != nil {
		_ = os.RemoveAll(dstPath) // Remove newly created dir because of error
		return err
	}
	c.Config.Generators[localName] = newGen // entry created

	return c.Save()
}

// Close if temporary removes src
func (c *Core) Close() {
	if c.IsTemporary {
		defer os.RemoveAll(c.ProjRoot)

	}
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
		var u *url.URL
		u, err = url.Parse(source)
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
			c.IsTemporary = true
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
