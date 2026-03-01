package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
)

// type DatabasePostgres struct {
// 	logger l.Logger
// 	db     *sql.DB
// }

// func NewDatabasePostgres(logger l.Logger) *DatabasePostgres {
// 	dbLogger := logger.WithFields("file_name", "database_postgres.go", "class_name", "DatabasePostgres")
// 	return &DatabasePostgres{
// 		logger: dbLogger,
// 	}
// }

type DatabasePostgres struct {
	logger         l.Logger
	db             *sql.DB
	migrationsPath string
}

func NewDatabasePostgres(logger l.Logger) *DatabasePostgres {
	return &DatabasePostgres{
		logger:         logger.WithFields("file_name", "database_postgres.go", "class_name", "DatabasePostgres"),
		migrationsPath: "db/migrations", // default
	}
}

func (d *DatabasePostgres) WithDB(db *sql.DB) *DatabasePostgres {
	d.db = db
	return d
}

func (d *DatabasePostgres) WithMigrationsPath(path string) *DatabasePostgres {
	d.migrationsPath = path
	return d
}

func (d *DatabasePostgres) Connect() (*sql.DB, error) {
	d.logger.Debug("Connect Postgres Database")

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	sslmode := "disable"
	timezone := "America/New_York"

	dataSourceName := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		user, pass, dbName, host, port, sslmode, timezone,
	)

	d.logger.Trace("Open connection to postgres db.", "dsn", dataSourceName)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	d.db = db
	return d.db, nil
}

func (d *DatabasePostgres) GetConnection() (*sql.DB, error) {
	d.logger.Debug("Get Postgres connection.")

	if d.db == nil {
		d.logger.Trace("No connection currently open, attempt to connect to db.")
		return d.Connect()
	}

	return d.db, nil
}

func (d *DatabasePostgres) MigrationStatus() (*model.MigrationStatus, error) {
	d.logger.Debug("Get migration status.")

	if _, err := d.GetConnection(); err != nil {
		return nil, err
	}

	m, err := d.getMigrate()
	if err != nil {
		d.logger.Error("Failed to get migrate for status.", "error", err.Error())
		return nil, err
	}

	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		d.logger.Error("Failed to get migration version.", "error", err.Error())
		return nil, err
	}

	absPath, err := filepath.Abs(d.migrationsPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	var latestVersion uint = 0
	total := 0
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}
		total++
		var version uint
		fmt.Sscanf(e.Name(), "%d", &version)
		if version > latestVersion {
			latestVersion = version
		}
	}

	status := &model.MigrationStatus{
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		Pending:        int(latestVersion) - int(currentVersion),
		Total:          total,
		Dirty:          dirty,
	}

	d.logger.Debug("Migration status retrieved.",
		"current_version", status.CurrentVersion,
		"latest_version", status.LatestVersion,
		"pending", status.Pending,
		"total", status.Total,
		"dirty", status.Dirty,
	)

	return status, nil
}

func (d *DatabasePostgres) MigrationDown() (*model.MigrationStatus, error) {
	d.logger.Debug("Remove postgres migrations.")

	status, err := d.MigrationStatus()
	if err != nil {
		return nil, err
	}

	if status.CurrentVersion == 0 {
		d.logger.Debug("No migrations applied, skipping down.")
		return status, nil
	}

	m, err := d.getMigrate()
	if err != nil {
		return nil, err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to migrate down.", "error", err.Error())
		return nil, err
	}

	return d.MigrationStatus()
}

func (d *DatabasePostgres) MigrationSteps(steps int8) (*model.MigrationStatus, error) {
	d.logger.Debug("Step through postgres migrations.", "steps", steps)

	status, err := d.MigrationStatus()
	if err != nil {
		return nil, err
	}

	if steps > 0 && int(steps) > status.Pending {
		return nil, model.NewBadMigrationCommandRequestError("requested %d step(s) up but only %d pending", steps, status.Pending)
	}

	if steps < 0 && int(-steps) > int(status.CurrentVersion) {
		return nil, model.NewBadMigrationCommandRequestError("requested %d step(s) down but only at version %d", -steps, status.CurrentVersion)
	}

	m, err := d.getMigrate()
	if err != nil {
		return nil, err
	}

	if err := m.Steps(int(steps)); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to migrate step.", "error", err.Error())
		return nil, err
	}

	return d.MigrationStatus()
}

func (d *DatabasePostgres) MigrationUp() (*model.MigrationStatus, error) {
	d.logger.Debug("Add postgres migrations.")

	status, err := d.MigrationStatus()
	if err != nil {
		return nil, err
	}

	if status.Pending == 0 {
		d.logger.Debug("No pending migrations, skipping up.")
		return status, nil
	}

	m, err := d.getMigrate()
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to migrate up.", "error", err.Error())
		return nil, err
	}

	return d.MigrationStatus()
}

func (d *DatabasePostgres) getMigrate() (*migrate.Migrate, error) {
	d.logger.Trace("Get postgres driver instance.")
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		d.logger.Error("Failed to get postgres instance while getting migrate.", "error", err.Error())
		return nil, err
	}

	absPath, err := filepath.Abs(d.migrationsPath)
	if err != nil {
		d.logger.Error("Failed to get absolute path for migrations.", "error", err.Error())
		return nil, err
	}

	d.logger.Trace("Get new migrate instance from postgres driver.", "path", absPath)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres", driver)

	return m, err
}

func (d *DatabasePostgres) Version() (*string, error) {
	d.logger.Debug("Get postgres version.")

	if _, err := d.GetConnection(); err != nil {
		d.logger.Error("Failed to get postgres connection while getting version.", "error", err.Error())
		return nil, err
	}

	var version string
	err := d.db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		d.logger.Error("Failed to select postgres version.", "error", err.Error())
		return nil, err
	}

	d.logger.Debug("Retrieved postgres version.", "version", version)
	return &version, nil
}
