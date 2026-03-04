package repository_test

import (
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ccthomas/board-game/internal/helper"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
	"github.com/ccthomas/board-game/internal/repository"
)

func newTestAbilityRepo(db *helper.TestDatabase) *repository.AbilityRepositoryPostgres {
	logger, _ := l.NewLoggerSlog()
	return repository.NewAbilityRepositoryPostgres(logger, db)
}

// --- Ability Upsert / GetByID ---

func TestAbility_Insert(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	expected := helper.CreateAbility(nil, nil, nil, nil, nil)
	err := repo.Upsert(expected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repo.GetByID(expected.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	helper.AssertAbility(t, &expected, result)
}

func TestAbility_Update(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	first := helper.CreateAbility(nil, nil, nil, nil, nil)
	err := repo.Upsert(first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := helper.CreateAbility(&first.ID, nil, nil, nil, nil)
	err = repo.Upsert(second)
	if err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}

	result, err := repo.GetByID(first.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	helper.AssertAbility(t, &second, result)
}

func TestAbility_GetByID_NotFound(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	result, err := repo.GetByID(uuid.NewString())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
}

func TestAbility_GetAll_FilterDeleted(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	nameA := "aaa_first"
	nameC := "ccc_third"
	nameB := "bbb_second_deleted"

	first := helper.CreateAbility(nil, &nameA, nil, nil, nil)
	first.DeletedAt = nil

	secondDeleted := helper.CreateAbility(nil, &nameB, nil, nil, nil)

	third := helper.CreateAbility(nil, &nameC, nil, nil, nil)
	third.DeletedAt = nil

	for _, a := range []model.Ability{first, secondDeleted, third} {
		if err := repo.Upsert(a); err != nil {
			t.Fatalf("unexpected error on upsert: %v", err)
		}
	}

	result, err := repo.GetAll()
	if err != nil {
		t.Fatalf("unexpected error on GetAll: %v", err)
	}

	data := *result
	if len(data) != 2 {
		t.Fatalf("expected 2 results, got %d", len(data))
	}

	helper.AssertAbility(t, &first, &data[0])
	helper.AssertAbility(t, &third, &data[1])
}

// --- Effects ---

func TestAbilityEffect_UpsertAndGet(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	ability := helper.CreateAbility(nil, nil, nil, nil, nil)
	if err := repo.Upsert(ability); err != nil {
		t.Fatalf("unexpected error upserting ability: %v", err)
	}

	effect := helper.CreateAbilityEffect(nil, &ability.ID, nil, nil, nil, nil)
	if err := repo.UpsertEffect(effect); err != nil {
		t.Fatalf("unexpected error upserting effect: %v", err)
	}

	results, err := repo.GetEffectsByAbilityID(ability.ID)
	if err != nil {
		t.Fatalf("unexpected error on GetEffectsByAbilityID: %v", err)
	}

	if len(*results) != 1 {
		t.Fatalf("expected 1 effect, got %d", len(*results))
	}

	helper.AssertAbilityEffect(t, &effect, &(*results)[0])
}

func TestAbilityEffect_UpdateEffect(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	ability := helper.CreateAbility(nil, nil, nil, nil, nil)
	if err := repo.Upsert(ability); err != nil {
		t.Fatalf("unexpected error upserting ability: %v", err)
	}

	first := helper.CreateAbilityEffect(nil, &ability.ID, nil, nil, nil, nil)
	if err := repo.UpsertEffect(first); err != nil {
		t.Fatalf("unexpected error upserting effect: %v", err)
	}

	// Re-upsert with same ID but different expression
	second := helper.CreateAbilityEffect(&first.ID, &ability.ID, nil, nil, nil, nil)
	if err := repo.UpsertEffect(second); err != nil {
		t.Fatalf("unexpected error on update upsert: %v", err)
	}

	results, err := repo.GetEffectsByAbilityID(ability.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	if len(*results) != 1 {
		t.Fatalf("expected 1 effect after update, got %d", len(*results))
	}

	helper.AssertAbilityEffect(t, &second, &(*results)[0])
}

func TestAbilityEffect_DeleteEffectsByAbilityID(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	ability := helper.CreateAbility(nil, nil, nil, nil, nil)
	if err := repo.Upsert(ability); err != nil {
		t.Fatalf("unexpected error upserting ability: %v", err)
	}

	effect1 := helper.CreateAbilityEffect(nil, &ability.ID, nil, nil, nil, nil)
	effect2 := helper.CreateAbilityEffect(nil, &ability.ID, nil, nil, nil, nil)

	for _, e := range []model.AbilityEffect{effect1, effect2} {
		if err := repo.UpsertEffect(e); err != nil {
			t.Fatalf("unexpected error upserting effect: %v", err)
		}
	}

	if err := repo.DeleteEffectsByAbilityID(ability.ID); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}

	results, err := repo.GetEffectsByAbilityID(ability.ID)
	if err != nil {
		t.Fatalf("unexpected error on get after delete: %v", err)
	}

	if len(*results) != 0 {
		t.Fatalf("expected 0 effects after delete, got %d", len(*results))
	}
}

func TestAbility_GetByID_LoadsEffects(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	ability := helper.CreateAbility(nil, nil, nil, nil, nil)
	if err := repo.Upsert(ability); err != nil {
		t.Fatalf("unexpected error upserting ability: %v", err)
	}

	effect := helper.CreateAbilityEffect(nil, &ability.ID, nil, nil, nil, nil)
	if err := repo.UpsertEffect(effect); err != nil {
		t.Fatalf("unexpected error upserting effect: %v", err)
	}

	result, err := repo.GetByID(ability.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Effects) != 1 {
		t.Fatalf("expected 1 effect loaded on ability, got %d", len(result.Effects))
	}

	helper.AssertAbilityEffect(t, &effect, &result.Effects[0])
}

// Goal is to test general error handling — same pattern as DamageType repo test.
func TestAbilityUpsert_Error(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.ability_effect")
	helper.CleanTable(t, db.DB, "game.ability")
	repo := newTestAbilityRepo(db)

	idGreaterThan36 := "1234567890123456789012345678901234567890" // 40 characters
	ability := helper.CreateAbility(&idGreaterThan36, nil, nil, nil, nil)

	err := repo.Upsert(ability)
	if err == nil {
		t.Fatal("expected error for oversized id, got nil")
	}
}
