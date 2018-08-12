package result

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

func (r *Repository) GetOne(id int) (domain.Result, error) {
	result := domain.Result{}
	return result, sqlx.Get(r.db, &result, "SELECT results.* FROM results WHERE id=$1", id)
}

func (r *Repository) Create(result domain.Result) (domain.Result, error) {
	queryResult, err := r.db.Exec(`INSERT INTO
		results (job_id, steps, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, result.JobID, result.Steps)
	if err != nil {
		return domain.Result{}, err
	}
	id, err := queryResult.LastInsertId()
	if err != nil {
		return domain.Result{}, err
	}

	return r.GetOne(int(id))
}
