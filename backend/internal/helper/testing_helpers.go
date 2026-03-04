// Package helper...TODO
package helper

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/ccthomas/board-game/internal/logger"
	l "github.com/ccthomas/board-game/internal/logger/mock"
	"github.com/ccthomas/board-game/internal/model"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/ccthomas/board-game/internal/database"
)

// TODO - generated during AI Iteration Loop. Left be as low priority item to clean up.
// Moving forward to focus on delivering initial working product
func AssertError() error {
	return &MockError{}
}

type MockError struct{}

func (m *MockError) Error() string { return "failed" }

func NewDummyMockedLogger(ctrl *gomock.Controller) *l.MockLogger {
	mockLogger := l.NewMockLogger(ctrl)
	mockLogger.EXPECT().WithFields(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Trace(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	return mockLogger
}

func NewTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() { container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	logger, err := logger.NewLoggerSlog()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if _, err := database.NewDatabasePostgres(logger).
		WithDB(db).
		WithMigrationsPath("../../db/migrations").
		MigrationUp(); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return &TestDatabase{DB: db}
}

// It doesn't need a real connection — it just needs to satisfy the *sql.DB type for the
// mock expectation since gomock.Any() will match it regardless
func NewTestSQLDB() *sql.DB {
	db, _ := sql.Open("postgres", "")
	return db
}

// func NewTestDatabase(t *testing.T) *sql.DB {
// 	t.Helper()
// 	ctx := context.Background()

// 	container, err := postgres.RunContainer(ctx,
// 		testcontainers.WithImage("postgres:16"),
// 		postgres.WithDatabase("testdb"),
// 		postgres.WithUsername("test"),
// 		postgres.WithPassword("test"),
// 		testcontainers.WithWaitStrategy(
// 			wait.ForListeningPort("5432/tcp"),
// 		),
// 	)
// 	if err != nil {
// 		t.Fatalf("failed to start postgres container: %v", err)
// 	}

// 	t.Cleanup(func() { container.Terminate(ctx) })

// 	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
// 	if err != nil {
// 		t.Fatalf("failed to get connection string: %v", err)
// 	}

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		t.Fatalf("failed to open db: %v", err)
// 	}

// 	logger, err := logger.NewLoggerSlog()
// 	if err != nil {
// 		t.Fatalf("failed to create logger: %v", err)
// 	}

// 	if err := database.NewDatabasePostgres(logger).
// 		WithMigrationsPath("../../db/migrations").
// 		MigrationUp(db); err != nil {
// 		t.Fatalf("failed to run migrations: %v", err)
// 	}

// 	return db
// }

func CleanTable(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()
	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}

func AssertTimeEqual(t *testing.T, field string, expected, actual *time.Time) {
	t.Helper()
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		t.Fatalf("%s: expected \"%v\" but got \"%v\"", field, expected, actual)
	}
	if expected.UnixNano() != actual.UnixNano() {
		t.Fatalf("%s: expected \"%d\" but got \"%d\"", field, expected.UnixNano(), actual.UnixNano())
	}
}

func AssertDamageType(t *testing.T, expected *model.DamageType, actual *model.DamageType) {
	if actual == nil {
		t.Fatal("expected damage type, got nil")
	}
	if actual.ID != expected.ID {
		t.Fatalf("expected id to equal expected: \"%s\" actual: \"%s\"", expected.ID, actual.ID)
	}
	if actual.Name != expected.Name {
		t.Fatalf("expected name to equal expected: \"%s\" actual: \"%s\"", expected.Name, actual.Name)
	}

	AssertTimeEqual(t, "created_at", expected.CreatedAt, actual.CreatedAt)
	AssertTimeEqual(t, "updated_at", expected.UpdatedAt, actual.UpdatedAt)
	AssertTimeEqual(t, "deleted_at", expected.DeletedAt, actual.DeletedAt)
}

func CreateDamageType(id *string, name *string, createdAt *time.Time, updatedAt *time.Time, deletedAt *time.Time) model.DamageType {
	defaultID := uuid.New().String()
	defaultName := uuid.New().String()
	now := time.Now()
	defaultCreatedAt := now.Add(-3 * time.Second)
	defaultUpdatedAt := now.Add(-2 * time.Second)
	defaultDeletedAt := now.Add(-1 * time.Second)

	dt := model.DamageType{
		ID:        defaultID,
		Name:      defaultName,
		CreatedAt: &defaultCreatedAt,
		UpdatedAt: &defaultUpdatedAt,
		DeletedAt: &defaultDeletedAt,
	}

	if id != nil {
		dt.ID = *id
	}
	if name != nil {
		dt.Name = *name
	}
	if createdAt != nil {
		dt.CreatedAt = createdAt
	}
	if updatedAt != nil {
		dt.UpdatedAt = updatedAt
	}
	if deletedAt != nil {
		dt.DeletedAt = deletedAt
	}

	return dt
}

// internal/testhelper/database.go

// TestDatabase wraps a *sql.DB to satisfy the database.Database interface for tests.
type TestDatabase struct {
	DB *sql.DB
}

func (t *TestDatabase) Connect() (*sql.DB, error) {
	return t.DB, nil
}

func (t *TestDatabase) GetConnection() (*sql.DB, error) {
	return t.DB, nil
}

func (t *TestDatabase) MigrationDown() (*model.MigrationStatus, error) {
	return nil, nil
}

func (t *TestDatabase) MigrationStatus() (*model.MigrationStatus, error) {
	return nil, nil
}

func (t *TestDatabase) MigrationSteps(steps int8) (*model.MigrationStatus, error) {
	return nil, nil
}

func (t *TestDatabase) MigrationUp() (*model.MigrationStatus, error) {
	return nil, nil
}

func (t *TestDatabase) Version() (*string, error) {
	return nil, nil
}

func SuppressLogs() {
	os.Setenv("LOG_TO_TERMINAL", "false")
	os.Setenv("LOG_FILE_NAME", os.DevNull)
}

func CreateAbility(id *string, name *string, pattern *model.TargetingPattern, rangeVal *int, deletedAt *time.Time) model.Ability {
	defaultID := uuid.New().String()
	defaultName := uuid.New().String()
	defaultPattern := model.TargetAdjacent
	defaultRange := 1
	now := time.Now()
	defaultCreatedAt := now.Add(-3 * time.Second)
	defaultUpdatedAt := now.Add(-2 * time.Second)
	defaultDeletedAt := now.Add(-1 * time.Second)

	a := model.Ability{
		ID:        defaultID,
		Name:      defaultName,
		Pattern:   defaultPattern,
		Range:     defaultRange,
		CreatedAt: &defaultCreatedAt,
		UpdatedAt: &defaultUpdatedAt,
		DeletedAt: &defaultDeletedAt,
		Effects:   []model.AbilityEffect{},
	}

	if id != nil {
		a.ID = *id
	}
	if name != nil {
		a.Name = *name
	}
	if pattern != nil {
		a.Pattern = *pattern
	}
	if rangeVal != nil {
		a.Range = *rangeVal
	}
	if deletedAt != nil {
		a.DeletedAt = deletedAt
	}

	return a
}

func AssertAbility(t *testing.T, expected *model.Ability, actual *model.Ability) {
	t.Helper()
	if actual == nil {
		t.Fatal("expected ability, got nil")
	}
	if actual.ID != expected.ID {
		t.Fatalf("expected id to equal expected: \"%s\" actual: \"%s\"", expected.ID, actual.ID)
	}
	if actual.Name != expected.Name {
		t.Fatalf("expected name to equal expected: \"%s\" actual: \"%s\"", expected.Name, actual.Name)
	}
	if actual.Pattern != expected.Pattern {
		t.Fatalf("expected pattern to equal expected: \"%s\" actual: \"%s\"", expected.Pattern, actual.Pattern)
	}
	if actual.Range != expected.Range {
		t.Fatalf("expected range to equal expected: \"%d\" actual: \"%d\"", expected.Range, actual.Range)
	}

	AssertTimeEqual(t, "created_at", expected.CreatedAt, actual.CreatedAt)
	AssertTimeEqual(t, "updated_at", expected.UpdatedAt, actual.UpdatedAt)
	AssertTimeEqual(t, "deleted_at", expected.DeletedAt, actual.DeletedAt)
}

func CreateAbilityEffect(id *string, abilityID *string, effectType *model.EffectType, alignment *model.TargetAlignment, damageTypeID *string, deletedAt *time.Time) model.AbilityEffect {
	defaultID := uuid.New().String()
	defaultAbilityID := uuid.New().String()
	defaultEffectType := model.Damage
	defaultAlignment := model.AlignAlly
	defaultDamageTypeID := uuid.New().String()
	defaultExpression := model.DiceExpression{NumDice: 2, DieType: model.D6, Modifier: 0}
	now := time.Now()
	defaultCreatedAt := now.Add(-3 * time.Second)
	defaultUpdatedAt := now.Add(-2 * time.Second)
	defaultDeletedAt := now.Add(-1 * time.Second)

	e := model.AbilityEffect{
		ID:           defaultID,
		AbilityID:    defaultAbilityID,
		EffectType:   defaultEffectType,
		Alignment:    defaultAlignment,
		DamageTypeID: defaultDamageTypeID,
		Expression:   defaultExpression,
		CreatedAt:    &defaultCreatedAt,
		UpdatedAt:    &defaultUpdatedAt,
		DeletedAt:    &defaultDeletedAt,
	}

	if id != nil {
		e.ID = *id
	}
	if abilityID != nil {
		e.AbilityID = *abilityID
	}
	if effectType != nil {
		e.EffectType = *effectType
	}
	if alignment != nil {
		e.Alignment = *alignment
	}
	if damageTypeID != nil {
		e.DamageTypeID = *damageTypeID
	}
	if deletedAt != nil {
		e.DeletedAt = deletedAt
	}

	return e
}

func AssertAbilityEffect(t *testing.T, expected *model.AbilityEffect, actual *model.AbilityEffect) {
	t.Helper()
	if actual == nil {
		t.Fatal("expected ability effect, got nil")
	}
	if actual.ID != expected.ID {
		t.Fatalf("expected id to equal expected: \"%s\" actual: \"%s\"", expected.ID, actual.ID)
	}
	if actual.AbilityID != expected.AbilityID {
		t.Fatalf("expected ability_id to equal expected: \"%s\" actual: \"%s\"", expected.AbilityID, actual.AbilityID)
	}
	if actual.EffectType != expected.EffectType {
		t.Fatalf("expected effect_type to equal expected: \"%s\" actual: \"%s\"", expected.EffectType, actual.EffectType)
	}
	if actual.Alignment != expected.Alignment {
		t.Fatalf("expected alignment to equal expected: \"%s\" actual: \"%s\"", expected.Alignment, actual.Alignment)
	}
	if actual.DamageTypeID != expected.DamageTypeID {
		t.Fatalf("expected damage_type_id to equal expected: \"%s\" actual: \"%s\"", expected.DamageTypeID, actual.DamageTypeID)
	}
	if actual.Expression.NumDice != expected.Expression.NumDice {
		t.Fatalf("expected expression.num_dice to equal expected: \"%d\" actual: \"%d\"", expected.Expression.NumDice, actual.Expression.NumDice)
	}
	if actual.Expression.DieType != expected.Expression.DieType {
		t.Fatalf("expected expression.die_type to equal expected: \"%d\" actual: \"%d\"", expected.Expression.DieType, actual.Expression.DieType)
	}
	if actual.Expression.Modifier != expected.Expression.Modifier {
		t.Fatalf("expected expression.modifier to equal expected: \"%d\" actual: \"%d\"", expected.Expression.Modifier, actual.Expression.Modifier)
	}

	AssertTimeEqual(t, "created_at", expected.CreatedAt, actual.CreatedAt)
	AssertTimeEqual(t, "updated_at", expected.UpdatedAt, actual.UpdatedAt)
	AssertTimeEqual(t, "deleted_at", expected.DeletedAt, actual.DeletedAt)
}
