package stepconfig

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

func (r *Repository) GetOne(id int) (domain.StepConfig, error) {
	stepConfig := domain.StepConfig{}
	return stepConfig, sqlx.Get(r.db, &stepConfig, "SELECT step_configs.* FROM step_configs WHERE id=$1", id)
}

func (r *Repository) Create(stepConfig domain.StepConfig) (domain.StepConfig, error) {
	result, err := r.db.Exec(`INSERT INTO
		step_configs (tag, variables, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)`, stepConfig.Tag, stepConfig.Variables)
	if err != nil {
		return domain.StepConfig{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.StepConfig{}, err
	}

	return r.GetOne(int(id))
}
