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

func TestDamageTypeGetAll_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	mockService.EXPECT().
		GetAll().
		Return(nil, h.AssertError())

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/damage-type", nil)
	rr := httptest.NewRecorder()

	controller.GetAll(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestDamageTypeGetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	expected := &[]m.DamageType{
		{ID: "fire", Name: "Fire"},
		{ID: "ice", Name: "Ice"},
	}

	mockService.EXPECT().
		GetAll().
		Return(expected, nil)

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/damage-type", nil)
	rr := httptest.NewRecorder()

	controller.GetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body []m.DamageType
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(body) != 2 {
		t.Fatalf("expected 2 results, got %d", len(body))
	}
}

// --- GetByID ---

func TestDamageTypeGetByID_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	mockService.EXPECT().
		GetByID("fire").
		Return(nil, h.AssertError())

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/damage-type/fire", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "fire"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestDamageTypeGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	mockService.EXPECT().
		GetByID("fire").
		Return(nil, nil)

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/damage-type/fire", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "fire"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestDamageTypeGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	expected := &m.DamageType{ID: "fire", Name: "Fire"}

	mockService.EXPECT().
		GetByID("fire").
		Return(expected, nil)

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/damage-type/fire", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "fire"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body m.DamageType
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.ID != "fire" {
		t.Fatalf("expected fire, got %s", body.ID)
	}
}

// --- Save ---

func TestDamageTypeSave_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPut, "/damage-type", bytes.NewBuffer([]byte("invalid json")))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDamageTypeSave_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	damageType := m.DamageType{Name: "Fire"}
	bodyBytes, _ := json.Marshal(damageType)

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(nil, h.AssertError())

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPut, "/damage-type", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestDamageTypeSave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	now := time.Now()
	damageType := m.DamageType{Name: "Fire"}
	bodyBytes, _ := json.Marshal(damageType)

	expected := &m.DamageType{ID: "generated-id", Name: "Fire", CreatedAt: &now, UpdatedAt: &now}

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(expected, nil)

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPut, "/damage-type", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body m.DamageType
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.ID != "generated-id" {
		t.Fatalf("expected generated-id, got %s", body.ID)
	}
}

// --- Delete ---

func TestDamageTypeDelete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	mockService.EXPECT().
		Delete("fire").
		Return(h.AssertError())

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodDelete, "/damage-type/fire", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "fire"})
	rr := httptest.NewRecorder()

	controller.Delete(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestDamageTypeDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	mockService.EXPECT().
		Delete("fire").
		Return(nil)

	controller := NewDamageTypeController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodDelete, "/damage-type/fire", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "fire"})
	rr := httptest.NewRecorder()

	controller.Delete(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}

// --- HandleSubrouter ---

func TestHandleDamageTypeSubrouter_WiresRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockDamageTypeService(ctrl)

	controller := NewDamageTypeController(mockLogger, mockService)
	router := mux.NewRouter()
	controller.HandleSubrouter(router)

	// GET /damage-type
	mockService.EXPECT().GetAll().Return(&[]m.DamageType{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/damage-type", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /damage-type expected %d, got %d", http.StatusOK, rr.Code)
	}

	// GET /damage-type/{id}
	mockService.EXPECT().GetByID("fire").Return(&m.DamageType{ID: "fire", Name: "Fire"}, nil)
	req = httptest.NewRequest(http.MethodGet, "/damage-type/fire", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /damage-type/fire expected %d, got %d", http.StatusOK, rr.Code)
	}

	// PUT /damage-type
	mockService.EXPECT().Save(gomock.Any()).Return(&m.DamageType{ID: "fire", Name: "Fire"}, nil)
	bodyBytes, _ := json.Marshal(m.DamageType{Name: "Fire"})
	req = httptest.NewRequest(http.MethodPut, "/damage-type", bytes.NewBuffer(bodyBytes))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("PUT /damage-type expected %d, got %d", http.StatusOK, rr.Code)
	}

	// DELETE /damage-type/{id}
	mockService.EXPECT().Delete("fire").Return(nil)
	req = httptest.NewRequest(http.MethodDelete, "/damage-type/fire", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("DELETE /damage-type/fire expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}
