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

type DamageTypeService interface {
	Delete(id string) error
	GetAll() (*[]model.DamageType, error)
	GetByID(id string) (*model.DamageType, error)
	Save(damageType model.DamageType) (*model.DamageType, error)
}

type DamageTypeServiceImpl struct {
	logger         l.Logger
	damageTypeRepo repository.DamageTypeRepository
}

func NewDamageTypeServiceImpl(logger l.Logger, damageTypeRepo repository.DamageTypeRepository) *DamageTypeServiceImpl {
	serviceLogger := logger.WithFields(
		"file_name", "damage_type_service.go",
		"class_name", "DamageTypeService",
	)

	return &DamageTypeServiceImpl{
		logger:         serviceLogger,
		damageTypeRepo: damageTypeRepo,
	}
}

func (s *DamageTypeServiceImpl) Delete(id string) error {
	s.logger.Debug("Delete damage type.", "id", id)

	s.logger.Trace("Getting damage type by id.", "id", id)
	existing, err := s.damageTypeRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get damage type by id.", "error", err.Error())
		return err
	}

	if existing == nil {
		s.logger.Warn("Damage type not found.", "id", id)
		return errors.New("damage type not found: " + id)
	}

	now := time.Now()
	existing.UpdatedAt = &now
	existing.DeletedAt = &now

	s.logger.Trace("Upserting damage type with deleted at.", "id", id)
	if err := s.damageTypeRepo.Upsert(*existing); err != nil {
		s.logger.Error("Failed to upsert damage type while deleting.", "error", err.Error())
		return err
	}

	s.logger.Debug("Damage type deleted successfully.", "id", id)
	return nil
}

func (s *DamageTypeServiceImpl) GetAll() (*[]model.DamageType, error) {
	s.logger.Debug("Get all damage types.")

	results, err := s.damageTypeRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all damage types.", "error", err.Error())
		return nil, err
	}

	return results, nil
}

func (s *DamageTypeServiceImpl) GetByID(id string) (*model.DamageType, error) {
	s.logger.Debug("Get damage type by id.", "id", id)

	result, err := s.damageTypeRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get damage type by id.", "error", err.Error())
		return nil, err
	}

	return result, nil
}

func (s *DamageTypeServiceImpl) Save(damageType model.DamageType) (*model.DamageType, error) {
	s.logger.Debug("Save damage type.", "id", damageType.ID)

	now := time.Now()

	// New record — generate id and set created_at
	if damageType.ID == "" {
		damageType.ID = uuid.New().String()
		damageType.CreatedAt = &now
		s.logger.Trace("Generated new id for damage type.", "id", damageType.ID)
	} else {
		existing, err := s.GetByID(damageType.ID)
		if err != nil {
			return nil, err
		}

		if !helper.AreTimesEqual(existing.CreatedAt, damageType.CreatedAt) ||
			!helper.AreTimesEqual(existing.UpdatedAt, damageType.UpdatedAt) ||
			!helper.AreTimesEqual(existing.DeletedAt, damageType.DeletedAt) {
			return nil, model.NewBadRequestChangingTimestampsError()
		}
	}

	// Always clear deleted_at and update updated_at
	damageType.DeletedAt = nil
	damageType.UpdatedAt = &now

	s.logger.Trace("Upserting damage type.", "id", damageType.ID)
	if err := s.damageTypeRepo.Upsert(damageType); err != nil {
		s.logger.Error("Failed to upsert damage type.", "error", err.Error())
		return nil, err
	}

	s.logger.Debug("Damage type saved successfully.", "id", damageType.ID)
	return &damageType, nil
}
