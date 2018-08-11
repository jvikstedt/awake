package domain

import (
	"time"

	"github.com/jvikstedt/awake"
)

type StepConfig struct {
	ID        int
	Tag       Tag
	Variables awake.Variables
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type StepConfigRepository interface {
	GetOne(int) (StepConfig, error)
	Create(StepConfig) (StepConfig, error)
}
