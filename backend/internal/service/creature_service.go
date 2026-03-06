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

type CreatureService interface {
	Delete(id string) error
	GetAll() (*[]model.Creature, error)
	GetByID(id string) (*model.Creature, error)
	Save(creature model.Creature) (*model.Creature, error)
}

type CreatureServiceImpl struct {
	logger       l.Logger
	creatureRepo repository.CreatureRepository
	abilitySvc   AbilityService
}

func NewCreatureServiceImpl(logger l.Logger, creatureRepo repository.CreatureRepository, abilitySvc AbilityService) *CreatureServiceImpl {
	return &CreatureServiceImpl{
		logger:       logger.WithFields("file_name", "creature_service.go", "class_name", "CreatureServiceImpl"),
		creatureRepo: creatureRepo,
		abilitySvc:   abilitySvc,
	}
}

func (s *CreatureServiceImpl) Delete(id string) error {
	s.logger.Debug("Delete creature.", "id", id)

	existing, err := s.creatureRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get creature by id.", "error", err.Error())
		return err
	}

	if existing == nil {
		s.logger.Warn("Creature not found.", "id", id)
		return errors.New("creature not found: " + id)
	}

	now := time.Now()
	existing.UpdatedAt = &now
	existing.DeletedAt = &now

	s.logger.Trace("Upserting creature with deleted_at.", "id", id)
	if err := s.creatureRepo.Upsert(*existing); err != nil {
		s.logger.Error("Failed to upsert creature while deleting.", "error", err.Error())
		return err
	}

	s.logger.Debug("Creature deleted successfully.", "id", id)
	return nil
}

func (s *CreatureServiceImpl) GetAll() (*[]model.Creature, error) {
	s.logger.Debug("Get all creatures.")

	creatures, err := s.creatureRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all creatures.", "error", err.Error())
		return nil, err
	}

	for i := range *creatures {
		if err := s.enrichSlots((*creatures)[i].Abilities); err != nil {
			return nil, err
		}
	}

	return creatures, nil
}

func (s *CreatureServiceImpl) GetByID(id string) (*model.Creature, error) {
	s.logger.Debug("Get creature by id.", "id", id)

	creature, err := s.creatureRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get creature by id.", "error", err.Error())
		return nil, err
	}

	if creature == nil {
		return nil, nil
	}

	if err := s.enrichSlots(creature.Abilities); err != nil {
		return nil, err
	}

	return creature, nil
}

func (s *CreatureServiceImpl) Save(creature model.Creature) (*model.Creature, error) {
	s.logger.Debug("Save creature.", "id", creature.ID)

	now := time.Now()

	// New record — generate id and set created_at
	if creature.ID == "" {
		creature.ID = uuid.New().String()
		creature.CreatedAt = &now
		s.logger.Trace("Generated new id for creature.", "id", creature.ID)
	} else {
		existing, err := s.creatureRepo.GetByID(creature.ID)
		if err != nil {
			return nil, err
		}

		if existing == nil {
			return nil, errors.New("creature not found: " + creature.ID)
		}

		if !helper.AreTimesEqual(existing.CreatedAt, creature.CreatedAt) ||
			!helper.AreTimesEqual(existing.UpdatedAt, creature.UpdatedAt) ||
			!helper.AreTimesEqual(existing.DeletedAt, creature.DeletedAt) {
			return nil, model.NewBadRequestChangingTimestampsError()
		}
	}

	// Validate and prepare slots
	resolvedAbilities := make(map[string]model.Ability, len(creature.Abilities))

	for i := range creature.Abilities {
		slot := &creature.Abilities[i]

		if slot.AbilityID == "" && slot.Ability.ID != "" {
			s.logger.Trace("Falling back to transient ability id for slot.", "ability_id", slot.Ability.ID)
			slot.AbilityID = slot.Ability.ID
		}

		if slot.AbilityID == "" {
			return nil, errors.New("ability slot is missing ability_id")
		}

		ability, err := s.abilitySvc.GetByID(slot.AbilityID)
		if err != nil {
			s.logger.Error("Failed to look up ability for slot.", "ability_id", slot.AbilityID, "error", err.Error())
			return nil, err
		}
		if ability == nil {
			return nil, errors.New("ability not found: " + slot.AbilityID)
		}

		resolvedAbilities[slot.AbilityID] = *ability

		slot.CreatureID = creature.ID
		if slot.CreatedAt == nil {
			slot.CreatedAt = &now
		}
		slot.UpdatedAt = &now
		slot.DeletedAt = nil
	}

	// Persist creature row
	creature.DeletedAt = nil
	creature.UpdatedAt = &now

	s.logger.Trace("Upserting creature.", "id", creature.ID)
	if err := s.creatureRepo.Upsert(creature); err != nil {
		s.logger.Error("Failed to upsert creature.", "error", err.Error())
		return nil, err
	}

	// Replace slots — delete old, insert new
	s.logger.Trace("Replacing ability slots for creature.", "id", creature.ID)
	if err := s.creatureRepo.DeleteSlotsByCreatureID(creature.ID); err != nil {
		s.logger.Error("Failed to delete old ability slots.", "error", err.Error())
		return nil, err
	}

	for _, slot := range creature.Abilities {
		if err := s.creatureRepo.UpsertSlot(slot); err != nil {
			s.logger.Error("Failed to upsert ability slot.", "error", err.Error(), "ability_id", slot.AbilityID)
			return nil, err
		}
	}

	// Re-attach resolved Ability structs to slots before returning
	for i := range creature.Abilities {
		slot := &creature.Abilities[i]
		if ability, ok := resolvedAbilities[slot.AbilityID]; ok {
			slot.Ability = ability
		}
	}

	s.logger.Debug("Creature saved successfully.", "id", creature.ID)
	return &creature, nil
}

// enrichSlots resolves the transient Ability field on each AbilitySlot.
func (s *CreatureServiceImpl) enrichSlots(slots []model.AbilitySlot) error {
	for i := range slots {
		slot := &slots[i]
		if slot.AbilityID == "" {
			continue
		}

		ability, err := s.abilitySvc.GetByID(slot.AbilityID)
		if err != nil {
			s.logger.Error("Failed to enrich ability on slot.", "ability_id", slot.AbilityID, "error", err.Error())
			return err
		}
		if ability != nil {
			slot.Ability = *ability
		}
	}
	return nil
}
