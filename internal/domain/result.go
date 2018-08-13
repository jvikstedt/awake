package domain

import "time"

type Result struct {
	ID        int        `json:"id"`
	JobID     int        `json:"jobID" db:"job_id"`
	Steps     *Steps     `json:"steps" db:"steps"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type ResultRepository interface {
	GetAll() ([]Result, error)
	GetOne(int) (Result, error)
	Create(Result) (Result, error)
}
