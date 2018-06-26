package task

import "github.com/jvikstedt/awake"

type StepConfig struct {
	Tag             `json:"tag"`
	awake.Variables `json:"variables"`
}
