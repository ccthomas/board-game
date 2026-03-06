package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
)

type CreatureRepositoryPostgres struct {
	logger   l.Logger
	database d.Database
}

func NewCreatureRepositoryPostgres(logger l.Logger, database d.Database) *CreatureRepositoryPostgres {
	return &CreatureRepositoryPostgres{
		logger:   logger.WithFields("file_name", "creature_repository_postgres.go", "class_name", "CreatureRepositoryPostgres"),
		database: database,
	}
}

func (r *CreatureRepositoryPostgres) getConnection() (*sql.DB, error) {
	db, err := r.database.GetConnection()
	if err != nil {
		r.logger.Error("Failed to get database connection.", "error", err.Error())
		return nil, err
	}
	return db, nil
}

// --- Creature ---

func (r *CreatureRepositoryPostgres) GetAll() (*[]model.Creature, error) {
	r.logger.Debug("Get all creatures.")

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, health_points, defence, initiative, movement, action_count,
		       created_at, updated_at, deleted_at
		FROM game.creature
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query creatures.", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	results := []model.Creature{}
	for rows.Next() {
		c, err := scanCreatureRow(rows)
		if err != nil {
			r.logger.Error("Failed to scan creature row.", "error", err.Error())
			return nil, err
		}

		slots, err := r.GetSlotsByCreatureID(c.ID)
		if err != nil {
			return nil, err
		}
		c.Abilities = *slots

		results = append(results, *c)
	}

	return &results, nil
}

func (r *CreatureRepositoryPostgres) GetByID(id string) (*model.Creature, error) {
	r.logger.Debug("Get creature by id.", "id", id)

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, health_points, defence, initiative, movement, action_count,
		       created_at, updated_at, deleted_at
		FROM game.creature
		WHERE id = $1
	`

	row := db.QueryRow(query, id)
	c, err := scanCreature(row)
	if err == sql.ErrNoRows {
		r.logger.Debug("No creature found.", "id", id)
		return nil, nil
	}
	if err != nil {
		r.logger.Error("Failed to scan creature.", "error", err.Error())
		return nil, err
	}

	slots, err := r.GetSlotsByCreatureID(c.ID)
	if err != nil {
		return nil, err
	}
	c.Abilities = *slots

	return c, nil
}

func (r *CreatureRepositoryPostgres) Upsert(c model.Creature) error {
	r.logger.Debug("Upsert creature.", "id", c.ID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.ID)
	if err != nil {
		r.logger.Error("Failed to parse creature id as UUID.", "id", c.ID, "error", err.Error())
		return err
	}

	query := `
		INSERT INTO game.creature (id, name, health_points, defence, initiative, movement, action_count,
		                           created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE
			SET name          = EXCLUDED.name,
			    health_points = EXCLUDED.health_points,
			    defence       = EXCLUDED.defence,
			    initiative    = EXCLUDED.initiative,
			    movement      = EXCLUDED.movement,
			    action_count  = EXCLUDED.action_count,
			    created_at    = EXCLUDED.created_at,
			    updated_at    = EXCLUDED.updated_at,
			    deleted_at    = EXCLUDED.deleted_at;
	`

	_, err = db.Exec(query,
		id, c.Name, c.HealthPoints, c.Defence, c.Initiative, c.Movement, c.ActionCount,
		c.CreatedAt, c.UpdatedAt, c.DeletedAt,
	)
	if err != nil {
		r.logger.Error("Failed to upsert creature.", "error", err.Error())
		return err
	}

	return nil
}

// --- Ability Slots ---

func (r *CreatureRepositoryPostgres) GetSlotsByCreatureID(creatureID string) (*[]model.AbilitySlot, error) {
	r.logger.Debug("Get ability slots by creature id.", "creature_id", creatureID)

	db, err := r.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT creature_id, ability_id, roll_threshold, created_at, updated_at, deleted_at
		FROM game.creature_ability_slot
		WHERE creature_id = $1
		  AND deleted_at IS NULL
		ORDER BY roll_threshold ASC
	`

	rows, err := db.Query(query, creatureID)
	if err != nil {
		r.logger.Error("Failed to query ability slots.", "error", err.Error(), "creature_id", creatureID)
		return nil, err
	}
	defer rows.Close()

	results := []model.AbilitySlot{}
	for rows.Next() {
		slot, err := scanAbilitySlotRow(rows)
		if err != nil {
			r.logger.Error("Failed to scan ability slot row.", "error", err.Error())
			return nil, err
		}
		results = append(results, *slot)
	}

	return &results, nil
}

func (r *CreatureRepositoryPostgres) UpsertSlot(slot model.AbilitySlot) error {
	r.logger.Debug("Upsert ability slot.", "creature_id", slot.CreatureID, "ability_id", slot.AbilityID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	creatureID, err := uuid.Parse(slot.CreatureID)
	if err != nil {
		r.logger.Error("Failed to parse creature_id as UUID.", "creature_id", slot.CreatureID, "error", err.Error())
		return err
	}

	abilityID, err := uuid.Parse(slot.AbilityID)
	if err != nil {
		r.logger.Error("Failed to parse ability_id as UUID.", "ability_id", slot.AbilityID, "error", err.Error())
		return err
	}

	query := `
		INSERT INTO game.creature_ability_slot (creature_id, ability_id, roll_threshold,
		                                        created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (creature_id, ability_id) DO UPDATE
			SET roll_threshold = EXCLUDED.roll_threshold,
			    created_at     = EXCLUDED.created_at,
			    updated_at     = EXCLUDED.updated_at,
			    deleted_at     = EXCLUDED.deleted_at;
	`

	_, err = db.Exec(query,
		creatureID, abilityID, slot.RollThreshold,
		slot.CreatedAt, slot.UpdatedAt, slot.DeletedAt,
	)
	if err != nil {
		r.logger.Error("Failed to upsert ability slot.", "error", err.Error(), "creature_id", slot.CreatureID, "ability_id", slot.AbilityID)
		return err
	}

	return nil
}

func (r *CreatureRepositoryPostgres) DeleteSlotsByCreatureID(creatureID string) error {
	r.logger.Debug("Delete ability slots by creature id.", "creature_id", creatureID)

	db, err := r.getConnection()
	if err != nil {
		return err
	}

	query := `DELETE FROM game.creature_ability_slot WHERE creature_id = $1`

	_, err = db.Exec(query, creatureID)
	if err != nil {
		r.logger.Error("Failed to delete ability slots.", "error", err.Error(), "creature_id", creatureID)
		return err
	}

	return nil
}

// --- helpers ---

func scanCreature(row *sql.Row) (*model.Creature, error) {
	var c model.Creature
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := row.Scan(
		&id, &c.Name, &c.HealthPoints, &c.Defence, &c.Initiative, &c.Movement, &c.ActionCount,
		&createdAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	c.ID = id.String()
	c.CreatedAt = &createdAt
	c.UpdatedAt = &updatedAt
	c.DeletedAt = deletedAt
	return &c, nil
}

func scanCreatureRow(rows *sql.Rows) (*model.Creature, error) {
	var c model.Creature
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := rows.Scan(
		&id, &c.Name, &c.HealthPoints, &c.Defence, &c.Initiative, &c.Movement, &c.ActionCount,
		&createdAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	c.ID = id.String()
	c.CreatedAt = &createdAt
	c.UpdatedAt = &updatedAt
	c.DeletedAt = deletedAt
	return &c, nil
}

func scanAbilitySlotRow(rows *sql.Rows) (*model.AbilitySlot, error) {
	var slot model.AbilitySlot
	var creatureID, abilityID uuid.UUID
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	err := rows.Scan(
		&creatureID, &abilityID, &slot.RollThreshold,
		&createdAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	slot.CreatureID = creatureID.String()
	slot.AbilityID = abilityID.String()
	slot.CreatedAt = &createdAt
	slot.UpdatedAt = &updatedAt
	slot.DeletedAt = deletedAt
	return &slot, nil
}
