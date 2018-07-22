package job

import (
	"github.com/jmoiron/sqlx"
	"github.com/jvikstedt/awake/internal/domain"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll() ([]*domain.Job, error) {
	jobs := []*domain.Job{}
	err := sqlx.Select(r.db, &jobs, `SELECT jobs.* FROM jobs WHERE jobs.deleted_at IS NULL`)
	return jobs, err
}
