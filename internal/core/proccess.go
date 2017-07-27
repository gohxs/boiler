package core

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

// ProcessFile as template using data
func ProcessFile(source, dest string, data map[string]interface{}) error {
	//fsrc, err := os.Open(source)
	tmpl, err := template.ParseFiles(source)
	if err != nil {
		return err
	}

	fdst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fdst.Close()

	tmpl.Execute(fdst, data)

	// Chmod to same as source
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	err = os.Chmod(dest, sourceinfo.Mode())

	return nil

}

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
	// Destination is dir
	err = os.MkdirAll(dstPath, sourceinfo.Mode())
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(srcPath)

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

func ProcessString(source string, data map[string]interface{}) (string, error) {
	t, err := template.New("t").Option("missingkey=zero").Parse(source)
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
	err = os.Chmod(dest, sourceinfo.Mode())

	return nil
}
