// Package database...TODO
package database

import (
	"database/sql"

	"github.com/ccthomas/board-game/internal/model"
)

type Database interface {
	Connect() (*sql.DB, error)
	GetConnection() (*sql.DB, error)
	MigrationDown() (*model.MigrationStatus, error)
	MigrationSteps(steps int8) (*model.MigrationStatus, error)
	MigrationUp() (*model.MigrationStatus, error)
	MigrationStatus() (*model.MigrationStatus, error)
	Version() (*string, error)
}
