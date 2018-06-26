package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/jvikstedt/awake"
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

	stepConfigs := []task.StepConfig{
		task.StepConfig{
			Tag: "builtin_http",
			Variables: awake.Variables{
				"url": awake.Variable{
					Type: "string",
					Val:  "https://www.google.fi/",
				},
			},
		},
		task.StepConfig{
			Tag: "builtin_equal",
			Variables: awake.Variables{
				"actual": awake.Variable{
					Type: "dynamic",
					Val:  "${0:code}",
				},
				"expected": awake.Variable{
					Type: "integer",
					Val:  200,
				},
			},
		},
	}

	job := task.Job{
		ID:          1,
		Cron:        "@every 5s",
		StepConfigs: stepConfigs,
	}

	scheduler := cron.New(logger)

	scheduler.AddEntry(cron.EntryID(job.ID), job.Cron, func(id cron.EntryID) {
		t := task.New(logger, job.StepConfigs)
		steps := t.Run()

		data, _ := json.MarshalIndent(steps, "", "  ")
		logger.Printf("%s\n", data)
	})

	scheduler.Start()
	defer scheduler.Stop()

	// task := task.New(logger, stepConfigs)

	// steps := task.Run()
	// data, _ := json.MarshalIndent(steps, "", "  ")
	// logger.Printf("%s\n", data)
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
