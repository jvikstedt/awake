package domain

import "github.com/jvikstedt/awake"

type StepConfig struct {
	ID              int    `json:"id"`
	DisplayName     string `json:"displayName"`
	Tag             `json:"tag"`
	awake.Variables `json:"variables"`
}
