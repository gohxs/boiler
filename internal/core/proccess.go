package core

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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

// ProcessFile as template using data
func ProcessFile(source, dest string, data interface{}) error {
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

// Go trough each files and convert .t or .boiler files
func ProcessDir(source, dest string, data interface{}) error {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	dir, err := os.Open(source)
	if err != nil {
		return err
	}
	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, v := range entries {
		// Ignore git or others?
		if v.Name() == ".git" {
			continue
		}
		var err error
		srcFile := filepath.Join(source, v.Name())
		dstFile := filepath.Join(dest, v.Name())
		if v.IsDir() {
			//log.Printf("[dir]  Processing: %s", srcFile)
			err = ProcessDir(srcFile, dstFile, data)
		} else if strings.HasSuffix(v.Name(), ".boiler") && !strings.Contains(source, ".boiler") { // boiler extension except in .boiler folder
			//log.Printf("[tmpl] Processing: %s - %s", srcFile, dstFile[:len(dstFile)-7])
			err = ProcessFile(srcFile, dstFile[:len(dstFile)-7], data)
		} else {
			//log.Printf("[file] Processing: %s", srcFile)
			err = CopyFile(srcFile, dstFile)
		}
		if err != nil {
			return err
		}

	}

	return nil
}
