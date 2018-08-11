package domain

import (
	"time"

	"github.com/jvikstedt/awake"
)

type StepConfig struct {
	ID        int             `json:"id"`
	Tag       Tag             `json:"tag"`
	Variables awake.Variables `json:"variables"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
}

type StepConfigRepository interface {
	GetOne(int) (StepConfig, error)
	Create(StepConfig) (StepConfig, error)
}
