package database

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jvikstedt/awake/internal/job"
	"github.com/jvikstedt/awake/internal/result"
	_ "github.com/mattn/go-sqlite3"
)

var schemas = []string{
	schema,
	job.Schema,
	result.Schema,
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
	_, err := db.Exec(strings.Join(schemas, ""))
	return err
}
