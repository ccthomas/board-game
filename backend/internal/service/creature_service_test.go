package service

import (
	"errors"
	"testing"
	"time"

	h "github.com/ccthomas/board-game/internal/helper"
	"github.com/ccthomas/board-game/internal/model"
	rm "github.com/ccthomas/board-game/internal/repository/mock"
	sm "github.com/ccthomas/board-game/internal/service/mock"

	"go.uber.org/mock/gomock"
)

// --- helpers ---

func newCreatureService(ctrl *gomock.Controller, repo *rm.MockCreatureRepository, abilitySvc *sm.MockAbilityService) *CreatureServiceImpl {
	mockLogger := h.NewDummyMockedLogger(ctrl)
	return NewCreatureServiceImpl(mockLogger, repo, abilitySvc)
}

func baseCreature() model.Creature {
	return model.Creature{
		Name:         "Goblin",
		HealthPoints: 10,
		Defence:      model.DiceExpression{NumDice: 1, DieType: model.D6, Modifier: 0},
		Initiative:   3,
		Movement:     2,
		ActionCount:  1,
	}
}

func baseSlot(creatureID string) model.AbilitySlot {
	return model.AbilitySlot{
		AbilityID:     "ability-1",
		CreatureID:    creatureID,
		RollThreshold: 10,
	}
}

// --- Delete ---

func TestCreatureService_Delete_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().GetByID("creature-1").Return(nil, errors.New("db error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	err := svc.Delete("creature-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreatureService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().GetByID("creature-1").Return(nil, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	err := svc.Delete("creature-1")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestCreatureService_Delete_UpsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	existing := &model.Creature{ID: "creature-1", Name: "Goblin"}
	mockRepo.EXPECT().GetByID("creature-1").Return(existing, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(errors.New("upsert error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	err := svc.Delete("creature-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreatureService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	existing := &model.Creature{ID: "creature-1", Name: "Goblin"}
	mockRepo.EXPECT().GetByID("creature-1").Return(existing, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(c model.Creature) error {
		if c.DeletedAt == nil {
			t.Fatal("expected deleted_at to be set")
		}
		if c.UpdatedAt == nil {
			t.Fatal("expected updated_at to be set")
		}
		return nil
	})

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	err := svc.Delete("creature-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// --- GetAll ---

func TestCreatureService_GetAll_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_GetAll_EnrichError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creatures := &[]model.Creature{
		{ID: "creature-1", Abilities: []model.AbilitySlot{baseSlot("creature-1")}},
	}
	mockRepo.EXPECT().GetAll().Return(creatures, nil)
	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(nil, errors.New("lookup error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_GetAll_Success_NoSlots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creatures := &[]model.Creature{
		{ID: "creature-1", Name: "Goblin", Abilities: []model.AbilitySlot{}},
	}
	mockRepo.EXPECT().GetAll().Return(creatures, nil)
	// AbilityService should NOT be called when there are no slots

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(*result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(*result))
	}
}

func TestCreatureService_GetAll_Success_EnrichesAbility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creatures := &[]model.Creature{
		{ID: "creature-1", Abilities: []model.AbilitySlot{baseSlot("creature-1")}},
	}
	ability := &model.Ability{ID: "ability-1", Name: "Slash"}

	mockRepo.EXPECT().GetAll().Return(creatures, nil)
	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(ability, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	slots := (*result)[0].Abilities
	if slots[0].Ability.ID != "ability-1" {
		t.Fatalf("expected ability to be enriched, got %v", slots[0].Ability)
	}
}

// --- GetByID ---

func TestCreatureService_GetByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().GetByID("creature-1").Return(nil, errors.New("db error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetByID("creature-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().GetByID("creature-1").Return(nil, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetByID("creature-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for not found")
	}
}

func TestCreatureService_GetByID_EnrichError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := &model.Creature{ID: "creature-1", Abilities: []model.AbilitySlot{baseSlot("creature-1")}}
	mockRepo.EXPECT().GetByID("creature-1").Return(creature, nil)
	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(nil, errors.New("lookup error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetByID("creature-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := &model.Creature{ID: "creature-1", Abilities: []model.AbilitySlot{baseSlot("creature-1")}}
	ability := &model.Ability{ID: "ability-1", Name: "Slash"}

	mockRepo.EXPECT().GetByID("creature-1").Return(creature, nil)
	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(ability, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.GetByID("creature-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "creature-1" {
		t.Fatalf("expected creature-1, got %s", result.ID)
	}
	if result.Abilities[0].Ability.ID != "ability-1" {
		t.Fatal("expected ability to be enriched")
	}
}

// --- Save (new record) ---

func TestCreatureService_Save_NewRecord_UpsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().Upsert(gomock.Any()).Return(errors.New("upsert error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(baseCreature())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_NewRecord_Success_NoSlots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(c model.Creature) error {
		if c.ID == "" {
			t.Fatal("expected id to be generated")
		}
		if c.CreatedAt == nil {
			t.Fatal("expected created_at to be set")
		}
		if c.UpdatedAt == nil {
			t.Fatal("expected updated_at to be set")
		}
		if c.DeletedAt != nil {
			t.Fatal("expected deleted_at to be nil")
		}
		return nil
	})
	mockRepo.EXPECT().DeleteSlotsByCreatureID(gomock.Any()).Return(nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(baseCreature())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.ID == "" {
		t.Fatal("expected id to be set on returned object")
	}
}

func TestCreatureService_Save_NewRecord_MissingAbilityID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{{AbilityID: "", RollThreshold: 10}}

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err == nil {
		t.Fatal("expected error for missing ability_id, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_NewRecord_FallsBackToTransientAbilityID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{
		{
			AbilityID:     "",
			Ability:       model.Ability{ID: "ability-1", Name: "Slash"},
			RollThreshold: 10,
		},
	}
	ability := &model.Ability{ID: "ability-1", Name: "Slash"}

	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(ability, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteSlotsByCreatureID(gomock.Any()).Return(nil)
	mockRepo.EXPECT().UpsertSlot(gomock.Any()).DoAndReturn(func(s model.AbilitySlot) error {
		if s.AbilityID != "ability-1" {
			t.Fatalf("expected ability_id to be set from transient ability, got %s", s.AbilityID)
		}
		return nil
	})

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Abilities[0].AbilityID != "ability-1" {
		t.Fatalf("expected ability_id to be set on returned slot, got %s", result.Abilities[0].AbilityID)
	}
}

func TestCreatureService_Save_NewRecord_AbilityLookupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{baseSlot("")}

	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(nil, errors.New("lookup error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_NewRecord_AbilityNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{baseSlot("")}

	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(nil, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err == nil {
		t.Fatal("expected error for missing ability, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_NewRecord_DeleteSlotsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{baseSlot("")}
	ability := &model.Ability{ID: "ability-1", Name: "Slash"}

	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(ability, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteSlotsByCreatureID(gomock.Any()).Return(errors.New("delete error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_NewRecord_UpsertSlotError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{baseSlot("")}
	ability := &model.Ability{ID: "ability-1", Name: "Slash"}

	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(ability, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteSlotsByCreatureID(gomock.Any()).Return(nil)
	mockRepo.EXPECT().UpsertSlot(gomock.Any()).Return(errors.New("upsert slot error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_NewRecord_WithSlot_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	creature := baseCreature()
	creature.Abilities = []model.AbilitySlot{baseSlot("")}
	ability := &model.Ability{ID: "ability-1", Name: "Slash"}

	mockAbilitySvc.EXPECT().GetByID("ability-1").Return(ability, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteSlotsByCreatureID(gomock.Any()).Return(nil)
	mockRepo.EXPECT().UpsertSlot(gomock.Any()).DoAndReturn(func(s model.AbilitySlot) error {
		if s.CreatureID == "" {
			t.Fatal("expected slot creature_id to be set")
		}
		if s.CreatedAt == nil {
			t.Fatal("expected slot created_at to be set")
		}
		if s.UpdatedAt == nil {
			t.Fatal("expected slot updated_at to be set")
		}
		if s.DeletedAt != nil {
			t.Fatal("expected slot deleted_at to be nil")
		}
		return nil
	})

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(creature)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Abilities[0].Ability.ID != "ability-1" {
		t.Fatal("expected ability to be re-attached on result")
	}
}

// --- Save (existing record) ---

func TestCreatureService_Save_ExistingRecord_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	now := time.Now()
	input := model.Creature{ID: "creature-1", Name: "Goblin", CreatedAt: &now, UpdatedAt: &now}

	mockRepo.EXPECT().GetByID("creature-1").Return(nil, errors.New("db error"))

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_ExistingRecord_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	now := time.Now()
	input := model.Creature{ID: "creature-1", Name: "Goblin", CreatedAt: &now, UpdatedAt: &now}

	mockRepo.EXPECT().GetByID("creature-1").Return(nil, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_ExistingRecord_ChangingCreatedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.Creature{ID: "creature-1", CreatedAt: &now, UpdatedAt: &now}
	input := model.Creature{ID: "creature-1", CreatedAt: &different, UpdatedAt: &now}

	mockRepo.EXPECT().GetByID("creature-1").Return(existing, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var badReqErr *model.BadRequestChangingTimestampsError
	if !errors.As(err, &badReqErr) {
		t.Fatalf("expected BadRequestChangingTimestampsError, got %T", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_ExistingRecord_ChangingUpdatedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.Creature{ID: "creature-1", CreatedAt: &now, UpdatedAt: &now}
	input := model.Creature{ID: "creature-1", CreatedAt: &now, UpdatedAt: &different}

	mockRepo.EXPECT().GetByID("creature-1").Return(existing, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var badReqErr *model.BadRequestChangingTimestampsError
	if !errors.As(err, &badReqErr) {
		t.Fatalf("expected BadRequestChangingTimestampsError, got %T", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_ExistingRecord_ChangingDeletedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.Creature{ID: "creature-1", CreatedAt: &now, UpdatedAt: &now, DeletedAt: nil}
	input := model.Creature{ID: "creature-1", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &different}

	mockRepo.EXPECT().GetByID("creature-1").Return(existing, nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var badReqErr *model.BadRequestChangingTimestampsError
	if !errors.As(err, &badReqErr) {
		t.Fatalf("expected BadRequestChangingTimestampsError, got %T", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestCreatureService_Save_ExistingRecord_ClearsDeletedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockCreatureRepository(ctrl)
	mockAbilitySvc := sm.NewMockAbilityService(ctrl)

	now := time.Now()
	existing := &model.Creature{ID: "creature-1", Name: "Goblin", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}
	input := model.Creature{ID: "creature-1", Name: "Goblin2", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}

	mockRepo.EXPECT().GetByID("creature-1").Return(existing, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(c model.Creature) error {
		if c.DeletedAt != nil {
			t.Fatal("expected deleted_at to be cleared")
		}
		return nil
	})
	mockRepo.EXPECT().DeleteSlotsByCreatureID("creature-1").Return(nil)

	svc := newCreatureService(ctrl, mockRepo, mockAbilitySvc)

	result, err := svc.Save(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.DeletedAt != nil {
		t.Fatal("expected deleted_at to be nil on returned object")
	}
}
