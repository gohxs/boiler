package boiler

import (
	"os"
	"path"
	"path/filepath"
)

// Find project root from start path, if not found will simple return start
func solveProjRoot(start string) (string, error) {
	cwd, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for f := cwd; f != ""; cwd, f = path.Split(cwd) {
		cwd = filepath.Clean(cwd)
		//log.Println("Cwd", cwd)
		boilerpath := filepath.Join(cwd, ".boiler")
		st, err := os.Stat(boilerpath)
		if err == nil && st.IsDir() { // ignore error
			return cwd, nil
		}
	}

	return start, nil
}
