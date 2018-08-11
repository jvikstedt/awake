package domain

import (
	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/plugin"
)

type Config struct {
	PerformerConfigs `json:"performerConfigs"`
}

type PerformerConfigs map[plugin.Tag]awake.Variables
