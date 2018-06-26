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
	pluginDir := filepath.Join(rootDir, "plugins")

	os.MkdirAll(rootDir, os.ModePerm)
	os.MkdirAll(pluginDir, os.ModePerm)

	return rootDir
}
