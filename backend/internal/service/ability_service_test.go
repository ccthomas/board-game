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

func newAbilityService(ctrl *gomock.Controller, repo *rm.MockAbilityRepository, dtSvc *sm.MockDamageTypeService) *AbilityServiceImpl {
	mockLogger := h.NewDummyMockedLogger(ctrl)
	return NewAbilityServiceImpl(mockLogger, repo, dtSvc)
}

func baseDamageEffect(abilityID string) model.AbilityEffect {
	return model.AbilityEffect{
		EffectType:   model.Damage,
		Alignment:    model.AlignEnemy,
		DamageTypeID: "fire-id",
		AbilityID:    abilityID,
	}
}

func baseAbility() model.Ability {
	return model.Ability{
		Name:    "Fireball",
		Pattern: model.TargetLine,
		Range:   3,
	}
}

// --- Delete ---

func TestAbilityService_Delete_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().GetByID("ability-1").Return(nil, errors.New("db error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	err := svc.Delete("ability-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAbilityService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().GetByID("ability-1").Return(nil, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	err := svc.Delete("ability-1")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestAbilityService_Delete_UpsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	existing := &model.Ability{ID: "ability-1", Name: "Fireball"}
	mockRepo.EXPECT().GetByID("ability-1").Return(existing, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(errors.New("upsert error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	err := svc.Delete("ability-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAbilityService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	existing := &model.Ability{ID: "ability-1", Name: "Fireball"}
	mockRepo.EXPECT().GetByID("ability-1").Return(existing, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(a model.Ability) error {
		if a.DeletedAt == nil {
			t.Fatal("expected deleted_at to be set")
		}
		if a.UpdatedAt == nil {
			t.Fatal("expected updated_at to be set")
		}
		return nil
	})

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	err := svc.Delete("ability-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// --- GetAll ---

func TestAbilityService_GetAll_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_GetAll_EnrichError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	abilities := &[]model.Ability{
		{ID: "ability-1", Effects: []model.AbilityEffect{baseDamageEffect("ability-1")}},
	}
	mockRepo.EXPECT().GetAll().Return(abilities, nil)
	mockDtSvc.EXPECT().GetByID("fire-id").Return(nil, errors.New("lookup error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_GetAll_Success_NoDamageEffects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	abilities := &[]model.Ability{
		{ID: "ability-1", Name: "Shield", Effects: []model.AbilityEffect{
			{EffectType: model.Defence, Alignment: model.AlignSelf},
		}},
	}
	mockRepo.EXPECT().GetAll().Return(abilities, nil)
	// DamageTypeService should NOT be called for non-damage effects

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(*result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(*result))
	}
}

func TestAbilityService_GetAll_Success_EnrichesDamageType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	abilities := &[]model.Ability{
		{ID: "ability-1", Effects: []model.AbilityEffect{baseDamageEffect("ability-1")}},
	}
	fireDt := &model.DamageType{ID: "fire-id", Name: "Fire"}

	mockRepo.EXPECT().GetAll().Return(abilities, nil)
	mockDtSvc.EXPECT().GetByID("fire-id").Return(fireDt, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	effects := (*result)[0].Effects
	if effects[0].DamageType.ID != "fire-id" {
		t.Fatalf("expected damage type to be enriched, got %v", effects[0].DamageType)
	}
}

// --- GetByID ---

func TestAbilityService_GetByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().GetByID("ability-1").Return(nil, errors.New("db error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetByID("ability-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().GetByID("ability-1").Return(nil, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetByID("ability-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for not found")
	}
}

func TestAbilityService_GetByID_EnrichError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := &model.Ability{ID: "ability-1", Effects: []model.AbilityEffect{baseDamageEffect("ability-1")}}
	mockRepo.EXPECT().GetByID("ability-1").Return(ability, nil)
	mockDtSvc.EXPECT().GetByID("fire-id").Return(nil, errors.New("lookup error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetByID("ability-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := &model.Ability{ID: "ability-1", Effects: []model.AbilityEffect{baseDamageEffect("ability-1")}}
	fireDt := &model.DamageType{ID: "fire-id", Name: "Fire"}

	mockRepo.EXPECT().GetByID("ability-1").Return(ability, nil)
	mockDtSvc.EXPECT().GetByID("fire-id").Return(fireDt, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.GetByID("ability-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "ability-1" {
		t.Fatalf("expected ability-1, got %s", result.ID)
	}
	if result.Effects[0].DamageType.ID != "fire-id" {
		t.Fatalf("expected damage type to be enriched")
	}
}

// --- Save (new record) ---

func TestAbilityService_Save_NewRecord_UpsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().Upsert(gomock.Any()).Return(errors.New("upsert error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(baseAbility())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_NewRecord_Success_NoEffects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(a model.Ability) error {
		if a.ID == "" {
			t.Fatal("expected id to be generated")
		}
		if a.CreatedAt == nil {
			t.Fatal("expected created_at to be set")
		}
		if a.UpdatedAt == nil {
			t.Fatal("expected updated_at to be set")
		}
		if a.DeletedAt != nil {
			t.Fatal("expected deleted_at to be nil")
		}
		return nil
	})
	mockRepo.EXPECT().DeleteEffectsByAbilityID(gomock.Any()).Return(nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(baseAbility())
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

func TestAbilityService_Save_NewRecord_WithDamageEffect_DamageTypeNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := baseAbility()
	ability.Effects = []model.AbilityEffect{baseDamageEffect("")}

	mockDtSvc.EXPECT().GetByID("fire-id").Return(nil, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(ability)
	if err == nil {
		t.Fatal("expected error for missing damage type, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_NewRecord_WithDamageEffect_DamageTypeLookupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := baseAbility()
	ability.Effects = []model.AbilityEffect{baseDamageEffect("")}

	mockDtSvc.EXPECT().GetByID("fire-id").Return(nil, errors.New("lookup error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(ability)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_NewRecord_WithDamageEffect_MissingDamageTypeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := baseAbility()
	ability.Effects = []model.AbilityEffect{
		{EffectType: model.Damage, Alignment: model.AlignEnemy, DamageTypeID: ""},
	}

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(ability)
	if err == nil {
		t.Fatal("expected error for missing damage_type_id, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_NewRecord_WithDamageEffect_DeleteEffectsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := baseAbility()
	ability.Effects = []model.AbilityEffect{baseDamageEffect("")}
	fireDt := &model.DamageType{ID: "fire-id", Name: "Fire"}

	mockDtSvc.EXPECT().GetByID("fire-id").Return(fireDt, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteEffectsByAbilityID(gomock.Any()).Return(errors.New("delete error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(ability)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_NewRecord_WithDamageEffect_UpsertEffectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := baseAbility()
	ability.Effects = []model.AbilityEffect{baseDamageEffect("")}
	fireDt := &model.DamageType{ID: "fire-id", Name: "Fire"}

	mockDtSvc.EXPECT().GetByID("fire-id").Return(fireDt, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteEffectsByAbilityID(gomock.Any()).Return(nil)
	mockRepo.EXPECT().UpsertEffect(gomock.Any()).Return(errors.New("effect upsert error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(ability)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_NewRecord_WithDamageEffect_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	ability := baseAbility()
	ability.Effects = []model.AbilityEffect{baseDamageEffect("")}
	fireDt := &model.DamageType{ID: "fire-id", Name: "Fire"}

	mockDtSvc.EXPECT().GetByID("fire-id").Return(fireDt, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteEffectsByAbilityID(gomock.Any()).Return(nil)
	mockRepo.EXPECT().UpsertEffect(gomock.Any()).DoAndReturn(func(e model.AbilityEffect) error {
		if e.ID == "" {
			t.Fatal("expected effect id to be generated")
		}
		if e.AbilityID == "" {
			t.Fatal("expected effect ability_id to be set")
		}
		if e.CreatedAt == nil {
			t.Fatal("expected effect created_at to be set")
		}
		if e.UpdatedAt == nil {
			t.Fatal("expected effect updated_at to be set")
		}
		if e.DeletedAt != nil {
			t.Fatal("expected effect deleted_at to be nil")
		}
		return nil
	})

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(ability)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Effects[0].DamageType.ID != "fire-id" {
		t.Fatal("expected damage type to be re-attached on result")
	}
}

// --- Save (existing record) ---

func TestAbilityService_Save_ExistingRecord_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	now := time.Now()
	input := model.Ability{ID: "ability-1", Name: "Fireball", CreatedAt: &now, UpdatedAt: &now}

	mockRepo.EXPECT().GetByID("ability-1").Return(nil, errors.New("db error"))

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_ExistingRecord_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	now := time.Now()
	input := model.Ability{ID: "ability-1", Name: "Fireball", CreatedAt: &now, UpdatedAt: &now}

	mockRepo.EXPECT().GetByID("ability-1").Return(nil, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

	result, err := svc.Save(input)
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestAbilityService_Save_ExistingRecord_ChangingCreatedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.Ability{ID: "ability-1", CreatedAt: &now, UpdatedAt: &now}
	input := model.Ability{ID: "ability-1", CreatedAt: &different, UpdatedAt: &now}

	mockRepo.EXPECT().GetByID("ability-1").Return(existing, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

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

func TestAbilityService_Save_ExistingRecord_ChangingUpdatedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.Ability{ID: "ability-1", CreatedAt: &now, UpdatedAt: &now}
	input := model.Ability{ID: "ability-1", CreatedAt: &now, UpdatedAt: &different}

	mockRepo.EXPECT().GetByID("ability-1").Return(existing, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

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

func TestAbilityService_Save_ExistingRecord_ChangingDeletedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.Ability{ID: "ability-1", CreatedAt: &now, UpdatedAt: &now, DeletedAt: nil}
	input := model.Ability{ID: "ability-1", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &different}

	mockRepo.EXPECT().GetByID("ability-1").Return(existing, nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

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

func TestAbilityService_Save_ExistingRecord_ClearsDeletedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := rm.NewMockAbilityRepository(ctrl)
	mockDtSvc := sm.NewMockDamageTypeService(ctrl)

	now := time.Now()
	existing := &model.Ability{ID: "ability-1", Name: "Fireball", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}
	input := model.Ability{ID: "ability-1", Name: "Fireball2", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}

	mockRepo.EXPECT().GetByID("ability-1").Return(existing, nil)
	mockRepo.EXPECT().Upsert(gomock.Any()).DoAndReturn(func(a model.Ability) error {
		if a.DeletedAt != nil {
			t.Fatal("expected deleted_at to be cleared")
		}
		return nil
	})
	mockRepo.EXPECT().DeleteEffectsByAbilityID("ability-1").Return(nil)

	svc := newAbilityService(ctrl, mockRepo, mockDtSvc)

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
