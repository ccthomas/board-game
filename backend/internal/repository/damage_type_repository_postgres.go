package repository

import (
	"database/sql"
	"time"

	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
)

type DamageTypeRepositoryPostgres struct {
	logger   l.Logger
	database d.Database
}

func NewDamageTypeRepositoryPostgres(logger l.Logger, database d.Database) *DamageTypeRepositoryPostgres {
	return &DamageTypeRepositoryPostgres{
		logger:   logger.WithFields("file_name", "damage_type_repository_postgres.go", "class_name", "DamageTypeRepositoryPostgres"),
		database: database,
	}
}

func (r *DamageTypeRepositoryPostgres) getConnection() (*sql.DB, error) {
	db, err := r.database.GetConnection()
	if err != nil {
		r.logger.Error("Failed to get database connection.", "error", err.Error())
		return nil, err
	}
	return db, nil
}

func (r *DamageTypeRepositoryPostgres) GetAll() (*[]model.DamageType, error) {
	r.logger.Debug("Get all damage types.")

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM game.damage_type
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query damage types.", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	results := []model.DamageType{}
	for rows.Next() {
		dt, err := scanDamageTypeRow(rows)
		if err != nil {
			r.logger.Error("Failed to scan damage type row.", "error", err.Error())
			return nil, err
		}
		results = append(results, *dt)
	}

	return &results, nil
}

func (r *DamageTypeRepositoryPostgres) GetByID(id string) (*model.DamageType, error) {
	r.logger.Debug("Get damage type by id.", "id", id)

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM game.damage_type
		WHERE id = $1
	`

	row := db.QueryRow(query, id)
	dt, err := scanDamageType(row)
	if err == sql.ErrNoRows {
		r.logger.Debug("No damage type found.", "id", id)
		return nil, nil
	}
	if err != nil {
		r.logger.Error("Failed to scan damage type.", "error", err.Error())
		return nil, err
	}

	return dt, nil
}

func (r *DamageTypeRepositoryPostgres) Upsert(d model.DamageType) error {
	r.logger.Debug("Upsert damage type.", "id", d.ID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO game.damage_type (id, name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
			SET name       = EXCLUDED.name,
			    created_at = EXCLUDED.created_at,
			    updated_at = EXCLUDED.updated_at,
			    deleted_at = EXCLUDED.deleted_at;
	`

	_, err = db.Exec(query, d.ID, d.Name, d.CreatedAt, d.UpdatedAt, d.DeletedAt)
	if err != nil {
		r.logger.Error("Failed to upsert damage type.", "error", err.Error())
		return err
	}

	return nil
}

// --- helpers ---

func scanDamageType(row *sql.Row) (*model.DamageType, error) {
	var dt model.DamageType
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := row.Scan(&dt.ID, &dt.Name, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	dt.CreatedAt = &createdAt
	dt.UpdatedAt = &updatedAt
	dt.DeletedAt = deletedAt
	return &dt, nil
}

func scanDamageTypeRow(rows *sql.Rows) (*model.DamageType, error) {
	var dt model.DamageType
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := rows.Scan(&dt.ID, &dt.Name, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	dt.CreatedAt = &createdAt
	dt.UpdatedAt = &updatedAt
	dt.DeletedAt = deletedAt
	return &dt, nil
}
