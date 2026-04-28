package persistence

import (
	"database/sql"
	_ "embed"
	"log"
)

//go:embed sql/schema.sql
var SchemaSQL string

// InitDB initializes the database schema
func InitDB(db *sql.DB) error {
	log.Printf("Creating Schema:\n%s\n", SchemaSQL)
	_, err := db.Exec(SchemaSQL)
	return err
}
