package repository

import (
	"github.com/ccthomas/board-game/internal/model"
)

type AbilityRepository interface {
	GetAll() (*[]model.Ability, error)
	GetByID(id string) (*model.Ability, error)
	Upsert(a model.Ability) error

	// Effects — scoped to a parent ability
	GetEffectsByAbilityID(abilityID string) (*[]model.AbilityEffect, error)
	UpsertEffect(e model.AbilityEffect) error
	DeleteEffectsByAbilityID(abilityID string) error
}
