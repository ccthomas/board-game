package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	h "github.com/ccthomas/board-game/internal/helper"
	m "github.com/ccthomas/board-game/internal/model"
	s "github.com/ccthomas/board-game/internal/service/mock"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

// --- GetAll ---

func TestAbilityGetAll_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	mockService.EXPECT().
		GetAll().
		Return(nil, h.AssertError())

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ability", nil)
	rr := httptest.NewRecorder()

	controller.GetAll(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestAbilityGetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	expected := &[]m.Ability{
		{ID: "ability-1", Name: "Fireball"},
		{ID: "ability-2", Name: "Shield"},
	}

	mockService.EXPECT().
		GetAll().
		Return(expected, nil)

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ability", nil)
	rr := httptest.NewRecorder()

	controller.GetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body []m.Ability
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(body) != 2 {
		t.Fatalf("expected 2 results, got %d", len(body))
	}
}

// --- GetByID ---

func TestAbilityGetByID_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	mockService.EXPECT().
		GetByID("ability-1").
		Return(nil, h.AssertError())

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ability/ability-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "ability-1"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestAbilityGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	mockService.EXPECT().
		GetByID("ability-1").
		Return(nil, nil)

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ability/ability-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "ability-1"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestAbilityGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	expected := &m.Ability{ID: "ability-1", Name: "Fireball"}

	mockService.EXPECT().
		GetByID("ability-1").
		Return(expected, nil)

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ability/ability-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "ability-1"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body m.Ability
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.ID != "ability-1" {
		t.Fatalf("expected ability-1, got %s", body.ID)
	}
}

// --- Save ---

func TestAbilitySave_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/ability", bytes.NewBuffer([]byte("invalid json")))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAbilitySave_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	ability := m.Ability{Name: "Fireball"}
	bodyBytes, _ := json.Marshal(ability)

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(nil, h.AssertError())

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/ability", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestAbilitySave_BadRequestChangingTimestamps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	ability := m.Ability{Name: "Fireball"}
	bodyBytes, _ := json.Marshal(ability)

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(nil, m.NewBadRequestChangingTimestampsError())

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/ability", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAbilitySave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	now := time.Now()
	ability := m.Ability{Name: "Fireball"}
	bodyBytes, _ := json.Marshal(ability)

	expected := &m.Ability{ID: "generated-id", Name: "Fireball", CreatedAt: &now, UpdatedAt: &now}

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(expected, nil)

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/ability", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body m.Ability
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.ID != "generated-id" {
		t.Fatalf("expected generated-id, got %s", body.ID)
	}
}

// --- Delete ---

func TestAbilityDelete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	mockService.EXPECT().
		Delete("ability-1").
		Return(h.AssertError())

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodDelete, "/ability/ability-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "ability-1"})
	rr := httptest.NewRecorder()

	controller.Delete(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestAbilityDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	mockService.EXPECT().
		Delete("ability-1").
		Return(nil)

	controller := NewAbilityController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodDelete, "/ability/ability-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "ability-1"})
	rr := httptest.NewRecorder()

	controller.Delete(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}

// --- HandleSubrouter ---

func TestHandleAbilitySubrouter_WiresRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockAbilityService(ctrl)

	controller := NewAbilityController(mockLogger, mockService)
	router := mux.NewRouter()
	controller.HandleSubrouter(router)

	// GET /ability
	mockService.EXPECT().GetAll().Return(&[]m.Ability{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/ability", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /ability expected %d, got %d", http.StatusOK, rr.Code)
	}

	// GET /ability/{id}
	mockService.EXPECT().GetByID("ability-1").Return(&m.Ability{ID: "ability-1", Name: "Fireball"}, nil)
	req = httptest.NewRequest(http.MethodGet, "/ability/ability-1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /ability/ability-1 expected %d, got %d", http.StatusOK, rr.Code)
	}

	// POST /ability
	mockService.EXPECT().Save(gomock.Any()).Return(&m.Ability{ID: "ability-1", Name: "Fireball"}, nil)
	bodyBytes, _ := json.Marshal(m.Ability{Name: "Fireball"})
	req = httptest.NewRequest(http.MethodPost, "/ability", bytes.NewBuffer(bodyBytes))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /ability expected %d, got %d", http.StatusOK, rr.Code)
	}

	// DELETE /ability/{id}
	mockService.EXPECT().Delete("ability-1").Return(nil)
	req = httptest.NewRequest(http.MethodDelete, "/ability/ability-1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("DELETE /ability/ability-1 expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}
