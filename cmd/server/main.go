package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/task"
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

	registerPerformers()

	stepConfigs := []task.StepConfig{
		task.StepConfig{
			Tag: "HTTP",
			Variables: awake.Variables{
				"url": awake.Variable{
					Type: "string",
					Val:  "https://www.google.fi/",
				},
			},
		},
		task.StepConfig{
			Tag: "EQUAL",
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

	task := task.New(logger, stepConfigs)

	steps := task.Run()
	data, _ := json.MarshalIndent(steps, "", "\t")
	logger.Printf("%s\n", data)
}
