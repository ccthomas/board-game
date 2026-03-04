package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
)

type AbilityRepositoryPostgres struct {
	logger   l.Logger
	database d.Database
}

func NewAbilityRepositoryPostgres(logger l.Logger, database d.Database) *AbilityRepositoryPostgres {
	return &AbilityRepositoryPostgres{
		logger:   logger.WithFields("file_name", "ability_repository_postgres.go", "class_name", "AbilityRepositoryPostgres"),
		database: database,
	}
}

func (r *AbilityRepositoryPostgres) getConnection() (*sql.DB, error) {
	db, err := r.database.GetConnection()
	if err != nil {
		r.logger.Error("Failed to get database connection.", "error", err.Error())
		return nil, err
	}
	return db, nil
}

// --- Ability ---

func (r *AbilityRepositoryPostgres) GetAll() (*[]model.Ability, error) {
	r.logger.Debug("Get all abilities.")

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, pattern, range, created_at, updated_at, deleted_at
		FROM game.ability
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query abilities.", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	results := []model.Ability{}
	for rows.Next() {
		a, err := scanAbilityRow(rows)
		if err != nil {
			r.logger.Error("Failed to scan ability row.", "error", err.Error())
			return nil, err
		}

		effects, err := r.GetEffectsByAbilityID(a.ID)
		if err != nil {
			return nil, err
		}
		a.Effects = *effects

		results = append(results, *a)
	}

	return &results, nil
}

func (r *AbilityRepositoryPostgres) GetByID(id string) (*model.Ability, error) {
	r.logger.Debug("Get ability by id.", "id", id)

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, pattern, range, created_at, updated_at, deleted_at
		FROM game.ability
		WHERE id = $1
	`

	row := db.QueryRow(query, id)
	a, err := scanAbility(row)
	if err == sql.ErrNoRows {
		r.logger.Debug("No ability found.", "id", id)
		return nil, nil
	}
	if err != nil {
		r.logger.Error("Failed to scan ability.", "error", err.Error())
		return nil, err
	}

	effects, err := r.GetEffectsByAbilityID(a.ID)
	if err != nil {
		return nil, err
	}
	a.Effects = *effects

	return a, nil
}

func (r *AbilityRepositoryPostgres) Upsert(a model.Ability) error {
	r.logger.Debug("Upsert ability.", "id", a.ID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	id, err := uuid.Parse(a.ID)
	if err != nil {
		r.logger.Error("Failed to parse ability id as UUID.", "id", a.ID, "error", err.Error())
		return err
	}

	query := `
		INSERT INTO game.ability (id, name, pattern, range, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
			SET name       = EXCLUDED.name,
			    pattern    = EXCLUDED.pattern,
			    range      = EXCLUDED.range,
			    created_at = EXCLUDED.created_at,
			    updated_at = EXCLUDED.updated_at,
			    deleted_at = EXCLUDED.deleted_at;
	`

	_, err = db.Exec(query, id, a.Name, a.Pattern, a.Range, a.CreatedAt, a.UpdatedAt, a.DeletedAt)
	if err != nil {
		r.logger.Error("Failed to upsert ability.", "error", err.Error())
		return err
	}

	return nil
}

// --- Effects ---

func (r *AbilityRepositoryPostgres) GetEffectsByAbilityID(abilityID string) (*[]model.AbilityEffect, error) {
	r.logger.Debug("Get effects by ability id.", "ability_id", abilityID)

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, ability_id, expression, effect_type, alignment, damage_type_id,
		       created_at, updated_at, deleted_at
		FROM game.ability_effect
		WHERE ability_id = $1
		ORDER BY created_at ASC
	`

	rows, err := db.Query(query, abilityID)
	if err != nil {
		r.logger.Error("Failed to query ability effects.", "error", err.Error(), "ability_id", abilityID)
		return nil, err
	}
	defer rows.Close()

	results := []model.AbilityEffect{}
	for rows.Next() {
		e, err := scanAbilityEffectRow(rows)
		if err != nil {
			r.logger.Error("Failed to scan ability effect row.", "error", err.Error())
			return nil, err
		}
		results = append(results, *e)
	}

	return &results, nil
}

func (r *AbilityRepositoryPostgres) UpsertEffect(e model.AbilityEffect) error {
	r.logger.Debug("Upsert ability effect.", "id", e.ID, "ability_id", e.AbilityID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	id, err := uuid.Parse(e.ID)
	if err != nil {
		r.logger.Error("Failed to parse effect id as UUID.", "id", e.ID, "error", err.Error())
		return err
	}

	abilityID, err := uuid.Parse(e.AbilityID)
	if err != nil {
		r.logger.Error("Failed to parse effect ability_id as UUID.", "ability_id", e.AbilityID, "error", err.Error())
		return err
	}

	damageTypeID := e.DamageTypeID
	if damageTypeID == "" && e.DamageType != nil {
		damageTypeID = e.DamageType.ID
	}

	query := `
		INSERT INTO game.ability_effect (id, ability_id, expression, effect_type, alignment, damage_type_id,
		                                 created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE
			SET ability_id     = EXCLUDED.ability_id,
			    expression     = EXCLUDED.expression,
			    effect_type    = EXCLUDED.effect_type,
			    alignment      = EXCLUDED.alignment,
			    damage_type_id = EXCLUDED.damage_type_id,
			    created_at     = EXCLUDED.created_at,
			    updated_at     = EXCLUDED.updated_at,
			    deleted_at     = EXCLUDED.deleted_at;
	`

	_, err = db.Exec(query,
		id, abilityID, e.Expression, e.EffectType, e.Alignment, damageTypeID,
		e.CreatedAt, e.UpdatedAt, e.DeletedAt,
	)
	if err != nil {
		r.logger.Error("Failed to upsert ability effect.", "error", err.Error(), "id", e.ID)
		return err
	}

	return nil
}

func (r *AbilityRepositoryPostgres) DeleteEffectsByAbilityID(abilityID string) error {
	r.logger.Debug("Delete effects by ability id.", "ability_id", abilityID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	query := `DELETE FROM game.ability_effect WHERE ability_id = $1`

	_, err = db.Exec(query, abilityID)
	if err != nil {
		r.logger.Error("Failed to delete ability effects.", "error", err.Error(), "ability_id", abilityID)
		return err
	}

	return nil
}

// --- helpers ---

func scanAbility(row *sql.Row) (*model.Ability, error) {
	var a model.Ability
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := row.Scan(&id, &a.Name, &a.Pattern, &a.Range, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	a.ID = id.String()
	a.CreatedAt = &createdAt
	a.UpdatedAt = &updatedAt
	a.DeletedAt = deletedAt
	return &a, nil
}

func scanAbilityRow(rows *sql.Rows) (*model.Ability, error) {
	var a model.Ability
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := rows.Scan(&id, &a.Name, &a.Pattern, &a.Range, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	a.ID = id.String()
	a.CreatedAt = &createdAt
	a.UpdatedAt = &updatedAt
	a.DeletedAt = deletedAt
	return &a, nil
}

func scanAbilityEffectRow(rows *sql.Rows) (*model.AbilityEffect, error) {
	var e model.AbilityEffect
	var id, abilityID uuid.UUID
	var damageTypeID uuid.NullUUID
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := rows.Scan(
		&id, &abilityID, &e.Expression, &e.EffectType, &e.Alignment, &damageTypeID,
		&createdAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	e.ID = id.String()
	e.AbilityID = abilityID.String()
	if damageTypeID.Valid {
		e.DamageTypeID = damageTypeID.UUID.String()
	}

	e.CreatedAt = &createdAt
	e.UpdatedAt = &updatedAt
	e.DeletedAt = deletedAt
	return &e, nil
}
