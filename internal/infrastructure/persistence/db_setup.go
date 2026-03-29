package persistence

import (
	"database/sql"
	_ "embed"
)

//go:embed sql/schema.sql
var SchemaSQL string

// InitDB initializes the database schema
func InitDB(db *sql.DB) error {
	_, err := db.Exec(SchemaSQL)
	return err
}
