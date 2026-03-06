package repository

import (
	"github.com/ccthomas/board-game/internal/model"
)

type CreatureRepository interface {
	GetAll() (*[]model.Creature, error)
	GetByID(id string) (*model.Creature, error)
	Upsert(c model.Creature) error

	// Ability slots — composite PK (creature_id, ability_id), no soft delete
	GetSlotsByCreatureID(creatureID string) (*[]model.AbilitySlot, error)
	UpsertSlot(slot model.AbilitySlot) error
	DeleteSlotsByCreatureID(creatureID string) error
}
