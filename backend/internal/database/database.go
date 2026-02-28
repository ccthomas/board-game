// Package database...TODO
package database

import "database/sql"

type Database interface {
	Connect() (*sql.DB, error)
	GetConnection() (*sql.DB, error)
	MigrationDown(db *sql.DB) error
	MigrationSteps(db *sql.DB, steps int8) error
	MigrationUp(db *sql.DB) error
	Version() (*string, error)
}
