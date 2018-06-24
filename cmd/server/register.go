package main

import (
	"fmt"
	"net/http"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/task"
)

func registerPerformers() {
	task.RegisterPerformer("EQUAL", Equal{})
	task.RegisterPerformer("HTTP", HTTP{})
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
