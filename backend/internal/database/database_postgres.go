package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	l "github.com/ccthomas/board-game/internal/logger"
)

type DatabasePostgres struct {
	logger l.Logger
	db     *sql.DB
}

func NewDatabasePostgres(logger l.Logger) *DatabasePostgres {
	dbLogger := logger.WithFields("file_name", "database_postgres.go", "class_name", "DatabasePostgres")
	return &DatabasePostgres{
		logger: dbLogger,
	}
}

func (d *DatabasePostgres) Connect() (*sql.DB, error) {
	d.logger.Debug("Connect Postgres Database")

	d.logger.Trace("Retrieve credentials from environment variables.")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// TODO make configurable
	sslmode := "disable" // local dev
	timezone := "America/New_York"

	dataSourceName := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		user,
		pass,
		dbName,
		host,
		port,
		sslmode,
		timezone,
	)

	d.logger.Trace("Open connection to postgres db.", "dsn", dataSourceName)
	db, err := sql.Open("postgres", dataSourceName)

	d.db = db
	return db, err
}

func (d *DatabasePostgres) GetConnection() (*sql.DB, error) {
	d.logger.Debug("Get Postgres connection.")

	if d.db == nil {
		d.logger.Trace("No connection currenly open, attempt to connect to db.")
		return d.Connect()
	}

	return d.db, nil
}

func (d *DatabasePostgres) MigrationDown(db *sql.DB) error {
	d.logger.Debug("Remove postgres migrations.")

	m, err := d.getMigrate(db)
	if err != nil {
		d.logger.Error("Failed to get migrate while migrating down", "error", err.Error())
		return err
	}

	d.logger.Trace("Run migrate down command.")
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to migrate down.", "error", err.Error())
		return err
	}

	return nil
}

func (d *DatabasePostgres) MigrationSteps(db *sql.DB, steps int8) error {
	d.logger.Debug("Step through postgres migrations.", "steps", steps)

	m, err := d.getMigrate(db)
	if err != nil {
		d.logger.Error("Failed to get migrate while migrating steps", "error", err.Error())
		return err
	}

	d.logger.Trace("Run migration step command.")
	if err := m.Steps(int(steps)); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to migrate step.", "error", err.Error())
		return err
	}

	return nil
}

func (d *DatabasePostgres) MigrationUp(db *sql.DB) error {
	d.logger.Debug("Add postgres migrations.")

	m, err := d.getMigrate(db)
	if err != nil {
		d.logger.Error("Failed to get migrate while migrating up", "error", err.Error())
		return err
	}

	d.logger.Trace("Run migration up command.")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to migrate up.", "error", err.Error())
		return err
	}

	return nil
}

func (d *DatabasePostgres) getMigrate(db *sql.DB) (*migrate.Migrate, error) {
	d.logger.Trace("Get postgres drive instance.")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		d.logger.Error("Failed to get postges instance while getting migrate.", "error", err.Error())
		return nil, err
	}

	absPath, err := filepath.Abs("db/migrations")
	if err != nil {
		d.logger.Error("Failed to get absolute path for migrations.", "error", err.Error())
		return nil, err
	}

	d.logger.Info("Resolved migrations path.", "abs_path", absPath)

	entries, readErr := os.ReadDir(absPath)
	if readErr != nil {
		d.logger.Error("Failed to read migrations directory.", "error", readErr.Error())
	} else {
		for _, e := range entries {
			d.logger.Info("Migration file found.", "file", e.Name())
		}
	}

	d.logger.Trace("Get new migrate instance from postgres driver.")
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath), // "file://db/migrations",
		"postgres", driver)

	return m, err
}

func (d *DatabasePostgres) Version() (*string, error) {
	d.logger.Debug("Get postgres version.")

	db, err := d.GetConnection()
	if err != nil {
		d.logger.Error("Failed to get postgres connection while getting version.", "error", err.Error())
		return nil, err
	}

	d.logger.Trace("Query postgres database for version")
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		d.logger.Error("Failed to select postgres version.", "error", err.Error())
		return nil, err
	}

	d.logger.Debug("Retrieved postgres version.", "verison", version)
	return &version, nil
}
