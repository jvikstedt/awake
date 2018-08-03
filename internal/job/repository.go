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

func (r *Repository) GetAll() ([]domain.Job, error) {
	jobs := []domain.Job{}
	err := sqlx.Select(r.db, &jobs, `SELECT jobs.* FROM jobs WHERE jobs.deleted_at IS NULL`)
	return jobs, err
}

func (r *Repository) GetOne(id int) (domain.Job, error) {
	job := domain.Job{}
	return job, sqlx.Get(r.db, &job, "SELECT jobs.* FROM jobs WHERE id=$1", id)
}

func (r *Repository) Update(id int, job domain.Job) (domain.Job, error) {
	_, err := r.db.Exec(`UPDATE jobs SET
		name = ?,
		active = ?,
		cron = ?,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, job.Name, job.Active, job.Cron, job.ID)
	if err != nil {
		return domain.Job{}, err
	}

	return r.GetOne(id)
}
