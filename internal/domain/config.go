package domain

import "github.com/jvikstedt/awake"

type Config struct {
	Jobs             []Job `json:"jobs"`
	PerformerConfigs `json:"performerConfigs"`
}

type PerformerConfigs map[Tag]awake.Variables
