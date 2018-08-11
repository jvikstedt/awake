package domain

import (
	"time"
)

type Job struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Active      bool         `json:"active"`
	Cron        string       `json:"cron"`
	StepConfigs *StepConfigs `json:"stepConfigs" db:"step_configs"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at" `
	UpdatedAt   time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time   `json:"deletedAt" db:"deleted_at"`
}

type JobRepository interface {
	GetAll() ([]Job, error)
	GetOne(int) (Job, error)
	Update(int, Job) (Job, error)
	Create(Job) (Job, error)
	Delete(int) (Job, error)
}
