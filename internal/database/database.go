package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/jvikstedt/awake/internal/job"
	_ "github.com/mattn/go-sqlite3"
)

var schemas = []string{
	schema,
	job.Schema,
}

func NewDB(driverName string, dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func EnsureTables(db *sqlx.DB) error {
	for _, s := range schemas {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}
