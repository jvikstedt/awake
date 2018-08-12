package domain

import (
	"github.com/jvikstedt/awake"
)

type StepResult struct {
	Variables awake.Variables `json:"variables"`
}
