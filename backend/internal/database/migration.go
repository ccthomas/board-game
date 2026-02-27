package database

import (
	"database/sql"
)

type Migration interface {
	Up(db *sql.DB) error
}
