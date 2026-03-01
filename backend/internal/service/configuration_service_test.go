package service

import (
	"errors"
	"testing"

	d "github.com/ccthomas/board-game/internal/database/mock"
	h "github.com/ccthomas/board-game/internal/helper"
	m "github.com/ccthomas/board-game/internal/model"

	"go.uber.org/mock/gomock"
)

func TestRunDatabaseMigration_Down_NoQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)

	mockDB.EXPECT().
		MigrationDown().
		Return(&m.MigrationStatus{}, nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationDown,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunDatabaseMigration_Down_NoQuantity_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)

	mockDB.EXPECT().
		MigrationDown().
		Return(nil, errors.New("migration failed"))

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationDown,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunDatabaseMigration_Down_WithQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	qty := int8(2)
	expectedSteps := int8(-2)

	mockDB.EXPECT().
		MigrationSteps(expectedSteps).
		Return(&m.MigrationStatus{}, nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command:  m.MigrationDown,
		Quantity: &qty,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunDatabaseMigration_Down_WithQuantity_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	qty := int8(2)
	expectedSteps := int8(-2)

	mockDB.EXPECT().
		MigrationSteps(expectedSteps).
		Return(nil, errors.New("migration failed"))

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command:  m.MigrationDown,
		Quantity: &qty,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunDatabaseMigration_Up_NoQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)

	mockDB.EXPECT().
		MigrationUp().
		Return(&m.MigrationStatus{}, nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationUp,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunDatabaseMigration_Up_NoQuantity_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)

	mockDB.EXPECT().
		MigrationUp().
		Return(nil, errors.New("migration failed"))

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationUp,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunDatabaseMigration_Up_WithQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	qty := int8(3)

	mockDB.EXPECT().
		MigrationSteps(qty).
		Return(&m.MigrationStatus{}, nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command:  m.MigrationUp,
		Quantity: &qty,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunDatabaseMigration_Up_WithQuantity_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	qty := int8(3)

	mockDB.EXPECT().
		MigrationSteps(qty).
		Return(nil, errors.New("migration failed"))

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command:  m.MigrationUp,
		Quantity: &qty,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunDatabaseMigration_UnknownCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: "invalid",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
