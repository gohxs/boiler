package core

import (
	"os"
	"path"
	"path/filepath"
)

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
