package service

import (
	"errors"
	"testing"

	d "github.com/ccthomas/board-game/internal/database/mock"
	h "github.com/ccthomas/board-game/internal/helper"
	"github.com/ccthomas/board-game/internal/model"

	"go.uber.org/mock/gomock"
)

func TestGetDatabaseHealth_VersionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := d.NewMockDatabase(ctrl)

	mockErr := errors.New("database unreachable")

	mockDB.EXPECT().
		Version().
		Return(nil, mockErr)

	service := NewHealthServiceImpl(h.NewDummyMockedLogger(ctrl), mockDB)

	result := service.GetDatabaseHealth()

	if result.Status != mockErr.Error() {
		t.Fatalf("expected status %s, got %s", mockErr.Error(), result.Status)
	}
}

func TestGetDatabaseHealth_VersionNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := d.NewMockDatabase(ctrl)

	mockDB.EXPECT().
		Version().
		Return(nil, nil)

	service := NewHealthServiceImpl(h.NewDummyMockedLogger(ctrl), mockDB)

	result := service.GetDatabaseHealth()

	if result.Status != "Healthy" {
		t.Fatalf("expected Healthy, got %s", result.Status)
	}
	if result.Version != "unknown" {
		t.Fatalf("expected unknown, got %s", result.Version)
	}
	if result.MigrationStatus != nil {
		t.Fatal("expected nil migration status")
	}
}

func TestGetDatabaseHealth_MigrationStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := d.NewMockDatabase(ctrl)
	version := "1.2.3"
	mockErr := errors.New("migration status failed")

	mockDB.EXPECT().
		Version().
		Return(&version, nil)

	mockDB.EXPECT().
		MigrationStatus().
		Return(nil, mockErr)

	service := NewHealthServiceImpl(h.NewDummyMockedLogger(ctrl), mockDB)

	result := service.GetDatabaseHealth()

	if result.Status != mockErr.Error() {
		t.Fatalf("expected status %s, got %s", mockErr.Error(), result.Status)
	}
}

func TestGetDatabaseHealth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := d.NewMockDatabase(ctrl)
	version := "1.2.3"
	migrationStatus := &model.MigrationStatus{
		CurrentVersion: 2,
		LatestVersion:  2,
		Pending:        0,
		Total:          2,
		Dirty:          false,
	}

	mockDB.EXPECT().
		Version().
		Return(&version, nil)

	mockDB.EXPECT().
		MigrationStatus().
		Return(migrationStatus, nil)

	service := NewHealthServiceImpl(h.NewDummyMockedLogger(ctrl), mockDB)

	result := service.GetDatabaseHealth()

	if result.Status != "Healthy" {
		t.Fatalf("expected Healthy, got %s", result.Status)
	}
	if result.Version != version {
		t.Fatalf("expected %s, got %s", version, result.Version)
	}
	if result.MigrationStatus == nil {
		t.Fatal("expected migration status, got nil")
	}
	if result.MigrationStatus.Pending != 0 {
		t.Fatalf("expected 0 pending, got %d", result.MigrationStatus.Pending)
	}
}
