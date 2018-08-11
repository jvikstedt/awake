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
	err := sqlx.Select(r.db, &jobs, `SELECT * FROM jobs WHERE jobs.deleted_at IS NULL`)
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
		step_config_ids = ?,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, job.Name, job.Active, job.Cron, job.StepConfigIDs, job.ID)
	if err != nil {
		return domain.Job{}, err
	}

	return r.GetOne(id)
}

func (r *Repository) Create(job domain.Job) (domain.Job, error) {
	result, err := r.db.Exec(`INSERT INTO
		jobs (name, active, cron, step_config_ids, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, job.Name, job.Active, job.Cron, job.StepConfigIDs)
	if err != nil {
		return domain.Job{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.Job{}, err
	}

	return r.GetOne(int(id))
}

func (r *Repository) Delete(id int) (domain.Job, error) {
	_, err := r.db.Exec(`UPDATE jobs SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, id)
	if err != nil {
		return domain.Job{}, err
	}

	return r.GetOne(id)
}
