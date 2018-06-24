package main

import (
	"os"
	"os/user"
	"path/filepath"
)

func getApplicationPath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	homeDir := usr.HomeDir

	rootDir := filepath.Join(homeDir, ".awake")

	os.MkdirAll(rootDir, os.ModePerm)
	return rootDir
}
