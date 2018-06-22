package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/task"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	task.RegisterPerformer("EQUAL", Equal{})
	task.RegisterPerformer("HTTP", HTTP{})

	steps := []task.Step{
		task.Step{
			Tag: "HTTP",
			Variables: awake.Variables{
				"url": awake.Variable{
					Type: "string",
					Val:  "https://www.google.fi/",
				},
			},
		},
		task.Step{
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

	task := task.New(logger, steps)

	task.Run()
}

type Equal struct{}

func (e Equal) Perform(scope awake.Scope) error {
	actual, _ := scope.ValueAsRaw("actual")
	expected, _ := scope.ValueAsRaw("expected")

	if actual != expected {
		return fmt.Errorf("Expected to be %v but got %v", expected, actual)
	}

	return nil
}

type HTTP struct{}

func (h HTTP) Perform(scope awake.Scope) error {
	url, _ := scope.ValueAsString("url")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	scope.SetReturnVariable("code", awake.Variable{
		Type: "integer",
		Val:  resp.StatusCode,
	})

	return nil
}
