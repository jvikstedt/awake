package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type StepConfigs []StepConfig

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

func (s *StepConfigs) Value() (driver.Value, error) {
	if s != nil {
		s, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(s), nil
	}
	return nil, nil
}

func (s *StepConfigs) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, s)
}
