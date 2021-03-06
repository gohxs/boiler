package boiler

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	funcMap template.FuncMap
)

func init() {

	funcMap = template.FuncMap{
		"Capitalize": strings.Title,
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
		"Base": func(name string) string {
			return filepath.Base(name)
		},
	}

}

// ProcessFile as template using data
func ProcessFile(source, dest string, data map[string]interface{}) error {
	//fsrc, err := os.Open(source)
	t, err := template.New(filepath.Base(source)).Funcs(funcMap).ParseFiles(source)
	if err != nil {
		log.Println("Has ERR", err)
		return err
	}

	_, err = os.Stat(dest)
	if !os.IsNotExist(err) {
		return os.ErrExist
	}

	fdst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fdst.Close()

	err = t.Execute(fdst, data)
	if err != nil {
		return err
	}

	// Chmod to same as source
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.Chmod(dest, sourceinfo.Mode())

}

// ProcessPath (either file or dir)
func ProcessPath(srcPath, dstPath string, data map[string]interface{}) error {
	sourceinfo, err := os.Stat(srcPath)
	if err != nil {
		log.Println("Err stat", err)
		return err
	}
	if !sourceinfo.IsDir() {
		if strings.HasSuffix(srcPath, BoilerExt) && !strings.Contains(dstPath, BoilerExt) { // boiler extension except in .boiler folder
			err = ProcessFile(srcPath, dstPath, data) // Process
		} else {
			err = CopyFile(srcPath, dstPath) // Simple copy
		}

		return err
	}
	// Source is dir so, Destination is dir
	err = os.MkdirAll(dstPath, sourceinfo.Mode())
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(srcPath)
	if err != nil {
		return err
	}

	for _, v := range entries {
		// Ignore git or others?
		if v.Name() == ".git" {
			continue
		}
		srcFile := filepath.Join(srcPath, v.Name())
		dstFile := filepath.Join(dstPath, v.Name())
		if v.Name() != BoilerExt && strings.HasSuffix(v.Name(), BoilerExt) && !strings.Contains(dstPath, BoilerExt) { // Except boiler path
			dstFile = dstFile[:len(dstFile)-7] // Remove boiler suffix // Remove suffix here?
		}
		err = ProcessPath(srcFile, dstFile, data) // Process ech file recursively
		if err != nil {
			return err
		}

	}

	return nil
}

//ProcessString passes a string through template with data
func ProcessString(source string, data map[string]interface{}) (string, error) {
	t, err := template.New("t").Option("missingkey=zero").Funcs(funcMap).Parse(source)
	if err != nil {
		return "", err
	}
	queryBuf := &bytes.Buffer{}
	t.Execute(queryBuf, data)

	return queryBuf.String(), nil
}

// CopyFile simple copy file from source to dest
func CopyFile(source, dest string) error {
	fsrc, err := os.Open(source)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	_, err = os.Stat(dest)
	if !os.IsNotExist(err) {
		return os.ErrExist
	}

	fdst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fdst.Close()

	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}
	// Chmod to same as source
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.Chmod(dest, sourceinfo.Mode())
}
