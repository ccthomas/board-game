package repository_test

import (
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ccthomas/board-game/internal/helper"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/repository"
)

func TestMain(m *testing.M) {
	helper.SuppressLogs()
	m.Run()
}

func newTestRepo(db *helper.TestDatabase) *repository.DamageTypeRepositoryPostgres {
	logger, _ := l.NewLoggerSlog()
	return repository.NewDamageTypeRepositoryPostgres(logger, db)
}

func Test_Insert(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.damage_type")
	repo := newTestRepo(db)

	expected := helper.CreateDamageType(nil, nil, nil, nil, nil)
	err := repo.Upsert(expected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repo.GetByID(expected.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	helper.AssertDamageType(t, &expected, result)
}

func Test_Update(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.damage_type")
	repo := newTestRepo(db)

	first := helper.CreateDamageType(nil, nil, nil, nil, nil)
	err := repo.Upsert(first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := helper.CreateDamageType(&first.ID, nil, nil, nil, nil)
	err = repo.Upsert(second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repo.GetByID(first.ID)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	helper.AssertDamageType(t, &second, result)
}

// Goal is to test general error handling.
// Not testing specific error 22001. Error 22001 is an easy known way to trigger an error.
func TestUpsert_Error(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.damage_type")
	repo := newTestRepo(db)

	idGreaterThan36 := "invalid-uuid" // 40 characters
	expected := helper.CreateDamageType(&idGreaterThan36, nil, nil, nil, nil)
	err := repo.Upsert(expected)
	if err == nil || err.Error() != "pq: invalid input syntax for type uuid: \"invalid-uuid\" (22P02)" {
		t.Fatalf("Expected error \"pq: invalid input syntax for type uuid: \"invalid-uuid\" (22P02)\": %v", err)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.damage_type")
	repo := newTestRepo(db)

	result, err := repo.GetByID(uuid.NewString())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Fatalf("Not result expected, but got: %v", result)
	}
}

func TestGetAll_FilterDeleted(t *testing.T) {
	db := helper.NewTestDatabase(t)
	helper.CleanTable(t, db.DB, "game.damage_type")
	repo := newTestRepo(db)

	nameA := "aaa_first"
	nameC := "ccc_third"
	nameB := "bbb_second_deleted"

	first := helper.CreateDamageType(nil, &nameA, nil, nil, nil)
	first.DeletedAt = nil

	secondDeleted := helper.CreateDamageType(nil, &nameB, nil, nil, nil)

	third := helper.CreateDamageType(nil, &nameC, nil, nil, nil)
	third.DeletedAt = nil

	err := repo.Upsert(first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Upsert(secondDeleted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Upsert(third)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repo.GetAll()
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}

	data := *result
	if len(data) != 2 {
		t.Fatalf("expected length of 2 for returned results: %v", len(data))
	}

	helper.AssertDamageType(t, &first, &data[0])
	helper.AssertDamageType(t, &third, &data[1])
}
