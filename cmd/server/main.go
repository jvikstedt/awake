package main

import (
	"fmt"

	global "github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/awake"
	"github.com/jvikstedt/awake/internal/task"
)

func main() {
	awake.RegisterPerformer("ASSERT_EQUAL", assertEqual{})

	steps := []awake.Step{
		awake.Step{
			Tag: "ASSERT_EQUAL",
		},
	}

	task := task.New(steps)

	task.Run()
}

type assertEqual struct {
}

func (ae assertEqual) Perform(scope global.Scope) error {
	fmt.Println("performing")
	return nil
}
