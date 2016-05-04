package common

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func CreatePathUnlessExist(relativePath string, perm os.FileMode) {
	if !FileExist(DefaultDataDir() + relativePath) {
		os.MkdirAll(DefaultDataDir()+relativePath, os.FileMode(perm))
	}
}

func ExpandHomePath(p string) (path string) {
	path = p
	sep := fmt.Sprintf("%s", os.PathSeparator)
	if len(p) > 1 && p[:1+len(sep)] == "~"+sep {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = strings.Replace(p, "~", dir, 1)
	}
	return
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func AbsolutePath(Datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(Datadir, filename)
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	var home string
	if usr, err := user.Current(); err == nil {
		home = usr.HomeDir
	} else {
		home = os.Getenv("HOME")
	}
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Oht")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Oht")
		} else {
			return filepath.Join(home, ".oht")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}
