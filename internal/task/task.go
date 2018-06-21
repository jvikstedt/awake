package task

import (
	"fmt"

	"github.com/jvikstedt/awake/internal/awake"
)

type Task struct {
	current   int
	stepPairs []*stepPair
}

type stepPair struct {
	step   awake.Step
	result awake.StepResult
	err    error
}

func New(steps []awake.Step) *Task {
	stepPairs := make([]*stepPair, len(steps))

	for i, step := range steps {
		stepPairs[i] = &stepPair{
			step:   step,
			result: awake.StepResult{},
		}
	}

	return &Task{
		current:   0,
		stepPairs: stepPairs,
	}
}

func (t *Task) Run() {
	for i, v := range t.stepPairs {
		t.current = i

		performer, ok := awake.FindPerformer(v.step.Tag)
		if !ok {
			fmt.Printf("Argh... performer not found %s\n", v.step.Tag)
			continue
		}

		err := performer.Perform(t)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Implements awake.Scope

func (t *Task) ValueAsRaw(name string) (interface{}, bool) {
	return nil, false
}

func (t *Task) ValueAsString(name string) (string, bool) {
	return "", false
}

func (t *Task) ValueAsInt(name string) (int, bool) {
	return 0, false
}

func (t *Task) ValueAsFloat(name string) (float64, bool) {
	return 0, false
}

func (t *Task) ValueAsBool(name string) (bool, bool) {
	return false, false
}
