package service

import (
	"errors"
	"testing"
	"time"

	h "github.com/ccthomas/board-game/internal/helper"
	"github.com/ccthomas/board-game/internal/model"
	r "github.com/ccthomas/board-game/internal/repository/mock"

	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	h.SuppressLogs()
	m.Run()
}

// --- Delete ---

func TestDamageTypeService_Delete_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	mockRepo.EXPECT().
		GetByID("fire").
		Return(nil, errors.New("db error"))

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	err := service.Delete("fire")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDamageTypeService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	mockRepo.EXPECT().
		GetByID("fire").
		Return(nil, nil)

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	err := service.Delete("fire")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestDamageTypeService_Delete_UpsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	existing := &model.DamageType{ID: "fire", Name: "Fire"}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(existing, nil)

	mockRepo.EXPECT().
		Upsert(gomock.Any()).
		Return(errors.New("upsert error"))

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	err := service.Delete("fire")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDamageTypeService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	existing := &model.DamageType{ID: "fire", Name: "Fire"}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(existing, nil)

	mockRepo.EXPECT().
		Upsert(gomock.Any()).
		DoAndReturn(func(d model.DamageType) error {
			if d.DeletedAt == nil {
				t.Fatal("expected deleted_at to be set")
			}
			if d.UpdatedAt == nil {
				t.Fatal("expected updated_at to be set")
			}
			return nil
		})

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	err := service.Delete("fire")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// --- GetAll ---

func TestDamageTypeService_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	mockRepo.EXPECT().
		GetAll().
		Return(nil, errors.New("db error"))

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.GetAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestDamageTypeService_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	expected := &[]model.DamageType{
		{ID: "fire", Name: "Fire"},
		{ID: "ice", Name: "Ice"},
	}

	mockRepo.EXPECT().
		GetAll().
		Return(expected, nil)

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(*result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(*result))
	}
}

// --- GetByID ---

func TestDamageTypeService_GetByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	mockRepo.EXPECT().
		GetByID("fire").
		Return(nil, errors.New("db error"))

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.GetByID("fire")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestDamageTypeService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	expected := &model.DamageType{ID: "fire", Name: "Fire"}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(expected, nil)

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.GetByID("fire")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "fire" {
		t.Fatalf("expected fire, got %s", result.ID)
	}
}

// --- Save ---

func TestDamageTypeService_Save_NewRecord_UpsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	mockRepo.EXPECT().
		Upsert(gomock.Any()).
		Return(errors.New("upsert error"))

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(model.DamageType{Name: "Fire"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestDamageTypeService_Save_NewRecord_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	mockRepo.EXPECT().
		Upsert(gomock.Any()).
		DoAndReturn(func(d model.DamageType) error {
			if d.ID == "" {
				t.Fatal("expected id to be generated")
			}
			if d.CreatedAt == nil {
				t.Fatal("expected created_at to be set")
			}
			if d.UpdatedAt == nil {
				t.Fatal("expected updated_at to be set")
			}
			if d.DeletedAt != nil {
				t.Fatal("expected deleted_at to be nil")
			}
			return nil
		})

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(model.DamageType{Name: "Fire"})
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

func TestDamageTypeService_Save_ExistingRecord_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	now := time.Now()
	input := model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &now}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(nil, errors.New("db error"))

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestDamageTypeService_Save_ExistingRecord_ChangingCreatedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &now}
	input := model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &different, UpdatedAt: &now}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(existing, nil)

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(input)
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

func TestDamageTypeService_Save_ExistingRecord_ChangingUpdatedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &now}
	input := model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &different}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(existing, nil)

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(input)
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

func TestDamageTypeService_Save_ExistingRecord_ChangingDeletedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	now := time.Now()
	different := now.Add(-1 * time.Hour)

	existing := &model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &now, DeletedAt: nil}
	input := model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &different}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(existing, nil)

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(input)
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

func TestDamageTypeService_Save_ExistingRecord_ClearsDeletedAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockRepo := r.NewMockDamageTypeRepository(ctrl)

	now := time.Now()
	// Existing record has DeletedAt set
	existing := &model.DamageType{ID: "fire", Name: "Fire", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}
	// Input matches existing timestamps exactly — passes the check
	input := model.DamageType{ID: "fire", Name: "Fire2", CreatedAt: &now, UpdatedAt: &now, DeletedAt: &now}

	mockRepo.EXPECT().
		GetByID("fire").
		Return(existing, nil)

	mockRepo.EXPECT().
		Upsert(gomock.Any()).
		DoAndReturn(func(d model.DamageType) error {
			if d.DeletedAt != nil {
				t.Fatal("expected deleted_at to be cleared after timestamp check")
			}
			return nil
		})

	service := NewDamageTypeServiceImpl(mockLogger, mockRepo)

	result, err := service.Save(input)
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
