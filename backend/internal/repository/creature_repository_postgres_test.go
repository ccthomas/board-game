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

func newTestCreatureRepo(db *helper.TestDatabase) *repository.CreatureRepositoryPostgres {
	logger, _ := l.NewLoggerSlog()
	return repository.NewCreatureRepositoryPostgres(logger, db)
}

func cleanCreatureTables(t *testing.T, db *helper.TestDatabase) {
	t.Helper()
	helper.CleanTable(t, db.DB, "game.creature_ability_slot")
	helper.CleanTable(t, db.DB, "game.creature")
}

// --- Creature Upsert / GetByID ---

func TestCreature_Insert(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	expected := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err := repo.Upsert(expected); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repo.GetByID(expected.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	helper.AssertCreature(t, &expected, result)
}

func TestCreature_Update(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	first := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err := repo.Upsert(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := helper.CreateCreature(&first.ID, nil, nil, nil, nil, nil, nil, nil, nil)
	if err := repo.Upsert(second); err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}

	result, err := repo.GetByID(first.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	helper.AssertCreature(t, &second, result)
}

func TestCreature_GetByID_NotFound(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	result, err := repo.GetByID(uuid.NewString())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
}

func TestCreature_GetAll_FilterDeleted(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	nameA := "aaa_first"
	nameB := "bbb_second_deleted"
	nameC := "ccc_third"

	first := helper.CreateCreature(nil, &nameA, nil, nil, nil, nil, nil, nil, nil)
	first.DeletedAt = nil

	secondDeleted := helper.CreateCreature(nil, &nameB, nil, nil, nil, nil, nil, nil, nil)

	third := helper.CreateCreature(nil, &nameC, nil, nil, nil, nil, nil, nil, nil)
	third.DeletedAt = nil

	for _, c := range []model.Creature{first, secondDeleted, third} {
		if err := repo.Upsert(c); err != nil {
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

	helper.AssertCreature(t, &first, &data[0])
	helper.AssertCreature(t, &third, &data[1])
}

func TestCreature_GetByID_LoadsSlots(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	creature := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	creature.DeletedAt = nil
	if err := repo.Upsert(creature); err != nil {
		t.Fatalf("unexpected error upserting creature: %v", err)
	}

	slot := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	slot.DeletedAt = nil
	if err := repo.UpsertSlot(slot); err != nil {
		t.Fatalf("unexpected error upserting slot: %v", err)
	}

	result, err := repo.GetByID(creature.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Abilities) != 1 {
		t.Fatalf("expected 1 slot loaded on creature, got %d", len(result.Abilities))
	}

	helper.AssertAbilitySlot(t, &slot, &result.Abilities[0])
}

func TestCreatureUpsert_Error(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	invalidID := "not-a-valid-uuid"
	creature := helper.CreateCreature(&invalidID, nil, nil, nil, nil, nil, nil, nil, nil)

	err := repo.Upsert(creature)
	if err == nil {
		t.Fatal("expected error for invalid uuid, got nil")
	}
}

// --- Ability Slots ---

func TestAbilitySlot_UpsertAndGet(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	creature := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	creature.DeletedAt = nil
	if err := repo.Upsert(creature); err != nil {
		t.Fatalf("unexpected error upserting creature: %v", err)
	}

	slot := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	slot.DeletedAt = nil
	if err := repo.UpsertSlot(slot); err != nil {
		t.Fatalf("unexpected error upserting slot: %v", err)
	}

	results, err := repo.GetSlotsByCreatureID(creature.ID)
	if err != nil {
		t.Fatalf("unexpected error on GetSlotsByCreatureID: %v", err)
	}

	if len(*results) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(*results))
	}

	helper.AssertAbilitySlot(t, &slot, &(*results)[0])
}

func TestAbilitySlot_Update(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	creature := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	creature.DeletedAt = nil
	if err := repo.Upsert(creature); err != nil {
		t.Fatalf("unexpected error upserting creature: %v", err)
	}

	first := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	first.DeletedAt = nil
	if err := repo.UpsertSlot(first); err != nil {
		t.Fatalf("unexpected error upserting slot: %v", err)
	}

	// Same composite PK, different roll_threshold
	newThreshold := first.RollThreshold + 5
	second := helper.CreateAbilitySlot(&creature.ID, &first.AbilityID, &newThreshold)
	second.DeletedAt = nil
	if err := repo.UpsertSlot(second); err != nil {
		t.Fatalf("unexpected error on update upsert: %v", err)
	}

	results, err := repo.GetSlotsByCreatureID(creature.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	if len(*results) != 1 {
		t.Fatalf("expected 1 slot after update, got %d", len(*results))
	}

	helper.AssertAbilitySlot(t, &second, &(*results)[0])
}

func TestAbilitySlot_FilterDeleted(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	creature := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	creature.DeletedAt = nil
	if err := repo.Upsert(creature); err != nil {
		t.Fatalf("unexpected error upserting creature: %v", err)
	}

	active := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	active.DeletedAt = nil

	deleted := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	// deleted.DeletedAt is set by default in CreateAbilitySlot

	for _, s := range []model.AbilitySlot{active, deleted} {
		if err := repo.UpsertSlot(s); err != nil {
			t.Fatalf("unexpected error upserting slot: %v", err)
		}
	}

	results, err := repo.GetSlotsByCreatureID(creature.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	if len(*results) != 1 {
		t.Fatalf("expected 1 active slot, got %d", len(*results))
	}

	helper.AssertAbilitySlot(t, &active, &(*results)[0])
}

func TestAbilitySlot_DeleteSlotsByCreatureID(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	creature := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	creature.DeletedAt = nil
	if err := repo.Upsert(creature); err != nil {
		t.Fatalf("unexpected error upserting creature: %v", err)
	}

	slot1 := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	slot1.DeletedAt = nil
	slot2 := helper.CreateAbilitySlot(&creature.ID, nil, nil)
	slot2.DeletedAt = nil

	for _, s := range []model.AbilitySlot{slot1, slot2} {
		if err := repo.UpsertSlot(s); err != nil {
			t.Fatalf("unexpected error upserting slot: %v", err)
		}
	}

	if err := repo.DeleteSlotsByCreatureID(creature.ID); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}

	results, err := repo.GetSlotsByCreatureID(creature.ID)
	if err != nil {
		t.Fatalf("unexpected error on get after delete: %v", err)
	}

	if len(*results) != 0 {
		t.Fatalf("expected 0 slots after delete, got %d", len(*results))
	}
}

func TestAbilitySlot_OrderedByRollThreshold(t *testing.T) {
	db := helper.NewTestDatabase(t)
	cleanCreatureTables(t, db)
	repo := newTestCreatureRepo(db)

	creature := helper.CreateCreature(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	creature.DeletedAt = nil
	if err := repo.Upsert(creature); err != nil {
		t.Fatalf("unexpected error upserting creature: %v", err)
	}

	high := 20
	mid := 10
	low := 1

	slotHigh := helper.CreateAbilitySlot(&creature.ID, nil, &high)
	slotHigh.DeletedAt = nil
	slotMid := helper.CreateAbilitySlot(&creature.ID, nil, &mid)
	slotMid.DeletedAt = nil
	slotLow := helper.CreateAbilitySlot(&creature.ID, nil, &low)
	slotLow.DeletedAt = nil

	for _, s := range []model.AbilitySlot{slotHigh, slotMid, slotLow} {
		if err := repo.UpsertSlot(s); err != nil {
			t.Fatalf("unexpected error upserting slot: %v", err)
		}
	}

	results, err := repo.GetSlotsByCreatureID(creature.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	data := *results
	if len(data) != 3 {
		t.Fatalf("expected 3 slots, got %d", len(data))
	}

	if data[0].RollThreshold != low {
		t.Fatalf("expected first slot roll_threshold %d, got %d", low, data[0].RollThreshold)
	}
	if data[1].RollThreshold != mid {
		t.Fatalf("expected second slot roll_threshold %d, got %d", mid, data[1].RollThreshold)
	}
	if data[2].RollThreshold != high {
		t.Fatalf("expected third slot roll_threshold %d, got %d", high, data[2].RollThreshold)
	}
}
