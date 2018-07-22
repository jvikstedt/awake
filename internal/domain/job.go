package domain

import "time"

type Job struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Active        bool         `json:"active"`
	Cron          string       `json:"cron"`
	StepConfigs   []StepConfig `json:"stepConfigs"`
	MailerEnabled bool         `json:"mailerEnabled"`
	CreatedAt     time.Time    `json:"createdAt" db:"created_at" `
	UpdatedAt     time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt     *time.Time   `json:"deletedAt" db:"deleted_at"`
}

type JobRepository interface {
	GetAll() ([]*Job, error)
}
