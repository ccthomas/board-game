package service

import (
	"database/sql"
	"errors"
	"testing"

	d "github.com/ccthomas/board-game/internal/database/mock"
	h "github.com/ccthomas/board-game/internal/helper"
	m "github.com/ccthomas/board-game/internal/model"

	"go.uber.org/mock/gomock"
)

func TestRunDatabaseMigration_GetConnectionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)

	mockDB.
		EXPECT().
		GetConnection().
		Return(nil, errors.New("connection failed"))

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationUp,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunDatabaseMigration_Down_NoQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	dbConn := &sql.DB{}

	mockDB.EXPECT().
		GetConnection().
		Return(dbConn, nil)

	mockDB.EXPECT().
		MigrationUp(dbConn).
		Return(nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationDown,
	})
	if err != nil {
		t.Fatal("No error should have been thrown")
	}
}

func TestRunDatabaseMigration_Down_WithQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	dbConn := &sql.DB{}
	qty := int8(2)
	expectedQty := int8(-2)

	mockDB.EXPECT().
		GetConnection().
		Return(dbConn, nil)

	mockDB.EXPECT().
		MigrationSteps(dbConn, qty).
		Return(nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command:  m.MigrationDown,
		Quantity: &expectedQty,
	})
	if err != nil {
		t.Fatal("No error should have been thrown")
	}
}

func TestRunDatabaseMigration_Up_NoQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	dbConn := &sql.DB{}

	mockDB.EXPECT().
		GetConnection().
		Return(dbConn, nil)

	mockDB.EXPECT().
		MigrationDown(dbConn).
		Return(nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: m.MigrationUp,
	})
	if err != nil {
		t.Fatal("No error should have been thrown")
	}
}

func TestRunDatabaseMigration_Up_WithQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	dbConn := &sql.DB{}
	qty := int8(3)
	expectedQty := int8(-3)

	mockDB.EXPECT().
		GetConnection().
		Return(dbConn, nil)

	mockDB.EXPECT().
		MigrationSteps(dbConn, qty*-1).
		Return(nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command:  m.MigrationUp,
		Quantity: &expectedQty,
	})
	if err != nil {
		t.Fatal("No error should have been thrown")
	}
}

func TestRunDatabaseMigration_UnknownCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockDB := d.NewMockDatabase(ctrl)
	dbConn := &sql.DB{}

	mockDB.EXPECT().
		GetConnection().
		Return(dbConn, nil)

	service := NewConfigurationServiceImpl(mockLogger, mockDB)

	err := service.RunDatabaseMigration(&m.MigrationConfiguration{
		Command: "invalid",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
