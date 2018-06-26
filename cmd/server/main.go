package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jvikstedt/awake/cron"
	"github.com/jvikstedt/awake/internal/task"
	"github.com/jvikstedt/awake/plugin"
)

func main() {
	appPath := getApplicationPath()

	// Setup logger
	f, err := os.OpenFile(filepath.Join(appPath, "awake.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)

	registerPerformers(logger, appPath)

	conf, err := loadConfig(logger, appPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scheduler := cron.New(logger)

	for _, j := range conf.Jobs {
		scheduler.AddEntry(cron.EntryID(j.ID), j.Cron, func(id cron.EntryID) {
			t := task.New(logger, j.StepConfigs)
			steps := t.Run()

			data, _ := json.MarshalIndent(steps, "", "  ")
			logger.Printf("%s\n", data)
		})
	}

	scheduler.Start()
	defer scheduler.Stop()
}

func registerPerformers(logger *log.Logger, appPath string) {
	plugin.BuiltinPerformers(func(performer task.Performer) {
		logger.Printf("Registering builtin performer %s\n", performer.Tag())
		task.RegisterPerformer(task.Tag(performer.Tag()), performer)
	})

	plugin.PluginPerformers(filepath.Join(appPath, "plugins"), func(performer task.Performer, err error) {
		logger.Printf("Registering plugin performer %s\n", performer.Tag())
		if err != nil {
			logger.Println(err)
			return
		}
		task.RegisterPerformer(task.Tag(performer.Tag()), performer)
	})
}

type config struct {
	Jobs []task.Job `json:"jobs"`
}

func loadConfig(logger *log.Logger, appPath string) (config, error) {
	data, err := ioutil.ReadFile(filepath.Join(appPath, "config.json"))
	if err != nil {
		return config{}, err
	}

	conf := config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return config{}, err
	}

	return conf, nil
}
