package domain

import "github.com/jvikstedt/awake"

type Config struct {
	PerformerConfigs `json:"performerConfigs"`
}

type PerformerConfigs map[Tag]awake.Variables
