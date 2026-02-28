package service

import (
	"errors"
	"testing"

	d "github.com/ccthomas/board-game/internal/database/mock"
	h "github.com/ccthomas/board-game/internal/helper"

	"go.uber.org/mock/gomock"
)

func TestGetDatabaseVersion_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := d.NewMockDatabase(ctrl)

	expected := "1.2.3"

	mockDB.
		EXPECT().
		Version().
		Return(&expected, nil)

	service := NewHealthServiceImpl(h.NewDummyMockedLogger(ctrl), mockDB)

	result := service.GetDatabaseVersion()

	if result != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestGetDatabaseVersion_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := d.NewMockDatabase(ctrl)

	mockErr := errors.New("database unreachable")

	mockDB.
		EXPECT().
		Version().
		Return(nil, mockErr)

	service := NewHealthServiceImpl(h.NewDummyMockedLogger(ctrl), mockDB)

	result := service.GetDatabaseVersion()

	if result != mockErr.Error() {
		t.Fatalf("expected error message %s, got %s", mockErr.Error(), result)
	}
}
