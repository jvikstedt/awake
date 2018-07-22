package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/jvikstedt/awake/internal/domain"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	appPath := getApplicationPath()

	// Setup logger
	f, err := os.OpenFile(filepath.Join(appPath, "awake.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)

	conf, err := loadConfig(logger, appPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app, err := newApp(logger, port, conf, appPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.registerPerformers()
	app.startServices()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		app.stopServices()
	}()

	app.wait()
}

func loadConfig(logger *log.Logger, appPath string) (domain.Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(appPath, "config.json"))
	if err != nil {
		return domain.Config{}, err
	}

	conf := domain.Config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return domain.Config{}, err
	}

	return conf, nil
}
