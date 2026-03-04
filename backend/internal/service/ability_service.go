package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ccthomas/board-game/internal/helper"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
	"github.com/ccthomas/board-game/internal/repository"
)

type AbilityService interface {
	Delete(id string) error
	GetAll() (*[]model.Ability, error)
	GetByID(id string) (*model.Ability, error)
	Save(ability model.Ability) (*model.Ability, error)
}

type AbilityServiceImpl struct {
	logger        l.Logger
	abilityRepo   repository.AbilityRepository
	damageTypeSvc DamageTypeService
}

func NewAbilityServiceImpl(logger l.Logger, abilityRepo repository.AbilityRepository, damageTypeSvc DamageTypeService) *AbilityServiceImpl {
	return &AbilityServiceImpl{
		logger:        logger.WithFields("file_name", "ability_service.go", "class_name", "AbilityServiceImpl"),
		abilityRepo:   abilityRepo,
		damageTypeSvc: damageTypeSvc,
	}
}

func (s *AbilityServiceImpl) Delete(id string) error {
	s.logger.Debug("Delete ability.", "id", id)

	existing, err := s.abilityRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get ability by id.", "error", err.Error())
		return err
	}

	if existing == nil {
		s.logger.Warn("Ability not found.", "id", id)
		return errors.New("ability not found: " + id)
	}

	now := time.Now()
	existing.UpdatedAt = &now
	existing.DeletedAt = &now

	s.logger.Trace("Upserting ability with deleted_at.", "id", id)
	if err := s.abilityRepo.Upsert(*existing); err != nil {
		s.logger.Error("Failed to upsert ability while deleting.", "error", err.Error())
		return err
	}

	s.logger.Debug("Ability deleted successfully.", "id", id)
	return nil
}

func (s *AbilityServiceImpl) GetAll() (*[]model.Ability, error) {
	s.logger.Debug("Get all abilities.")

	abilities, err := s.abilityRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all abilities.", "error", err.Error())
		return nil, err
	}

	for i := range *abilities {
		if err := s.enrichEffects((*abilities)[i].Effects); err != nil {
			return nil, err
		}
	}

	return abilities, nil
}

func (s *AbilityServiceImpl) GetByID(id string) (*model.Ability, error) {
	s.logger.Debug("Get ability by id.", "id", id)

	ability, err := s.abilityRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get ability by id.", "error", err.Error())
		return nil, err
	}

	if ability == nil {
		return nil, nil
	}

	if err := s.enrichEffects(ability.Effects); err != nil {
		return nil, err
	}

	return ability, nil
}

func (s *AbilityServiceImpl) Save(ability model.Ability) (*model.Ability, error) {
	s.logger.Debug("Save ability.", "id", ability.ID)

	now := time.Now()

	// New record — generate id and set created_at
	if ability.ID == "" {
		ability.ID = uuid.New().String()
		ability.CreatedAt = &now
		s.logger.Trace("Generated new id for ability.", "id", ability.ID)
	} else {
		existing, err := s.abilityRepo.GetByID(ability.ID)
		if err != nil {
			return nil, err
		}

		if existing == nil {
			return nil, errors.New("ability not found: " + ability.ID)
		}

		if !helper.AreTimesEqual(existing.CreatedAt, ability.CreatedAt) ||
			!helper.AreTimesEqual(existing.UpdatedAt, ability.UpdatedAt) ||
			!helper.AreTimesEqual(existing.DeletedAt, ability.DeletedAt) {
			return nil, model.NewBadRequestChangingTimestampsError()
		}
	}

	// Validate and prepare effects
	resolvedDamageTypes := make(map[string]model.DamageType, len(ability.Effects))

	for i := range ability.Effects {
		effect := &ability.Effects[i]

		// Validate damage effects have a damage type
		if effect.EffectType == model.Damage {
			if effect.DamageTypeID == "" {
				return nil, errors.New("damage effect is missing damage_type_id")
			}

			dt, err := s.damageTypeSvc.GetByID(effect.DamageTypeID)
			if err != nil {
				s.logger.Error("Failed to look up damage type for effect.", "damage_type_id", effect.DamageTypeID, "error", err.Error())
				return nil, err
			}
			if dt == nil {
				return nil, errors.New("damage type not found: " + effect.DamageTypeID)
			}

			resolvedDamageTypes[effect.DamageTypeID] = *dt
		}

		// Stamp effect ids and timestamps
		effect.AbilityID = ability.ID
		if effect.ID == "" {
			effect.ID = uuid.New().String()
			effect.CreatedAt = &now
			s.logger.Trace("Generated new id for ability effect.", "id", effect.ID)
		}
		effect.UpdatedAt = &now
		effect.DeletedAt = nil
	}

	// Persist ability row
	ability.DeletedAt = nil
	ability.UpdatedAt = &now

	s.logger.Trace("Upserting ability.", "id", ability.ID)
	if err := s.abilityRepo.Upsert(ability); err != nil {
		s.logger.Error("Failed to upsert ability.", "error", err.Error())
		return nil, err
	}

	// Replace effects — delete old, insert new
	s.logger.Trace("Replacing effects for ability.", "id", ability.ID)
	if err := s.abilityRepo.DeleteEffectsByAbilityID(ability.ID); err != nil {
		s.logger.Error("Failed to delete old ability effects.", "error", err.Error())
		return nil, err
	}

	for _, effect := range ability.Effects {
		if err := s.abilityRepo.UpsertEffect(effect); err != nil {
			s.logger.Error("Failed to upsert ability effect.", "error", err.Error(), "id", effect.ID)
			return nil, err
		}
	}

	// Re-attach resolved DamageType structs to effects before returning
	for i := range ability.Effects {
		effect := &ability.Effects[i]
		if effect.EffectType == model.Damage {
			if dt, ok := resolvedDamageTypes[effect.DamageTypeID]; ok {
				effect.DamageType = &dt
			}
		}
	}

	s.logger.Debug("Ability saved successfully.", "id", ability.ID)
	return &ability, nil
}

// enrichEffects resolves the DamageType transient field for any Damage effects.
func (s *AbilityServiceImpl) enrichEffects(effects []model.AbilityEffect) error {
	for i := range effects {
		effect := &effects[i]
		if effect.EffectType == model.Damage && effect.DamageTypeID != "" {
			dt, err := s.damageTypeSvc.GetByID(effect.DamageTypeID)
			if err != nil {
				s.logger.Error("Failed to enrich damage type on effect.", "damage_type_id", effect.DamageTypeID, "error", err.Error())
				return err
			}
			if dt != nil {
				effect.DamageType = dt
			}
		}
	}
	return nil
}
