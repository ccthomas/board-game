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

func TestCreatureGetAll_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	mockService.EXPECT().
		GetAll().
		Return(nil, h.AssertError())

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/creature", nil)
	rr := httptest.NewRecorder()

	controller.GetAll(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestCreatureGetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	expected := &[]m.Creature{
		{ID: "creature-1", Name: "Goblin"},
		{ID: "creature-2", Name: "Orc"},
	}

	mockService.EXPECT().
		GetAll().
		Return(expected, nil)

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/creature", nil)
	rr := httptest.NewRecorder()

	controller.GetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body []m.Creature
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(body) != 2 {
		t.Fatalf("expected 2 results, got %d", len(body))
	}
}

// --- GetByID ---

func TestCreatureGetByID_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	mockService.EXPECT().
		GetByID("creature-1").
		Return(nil, h.AssertError())

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/creature/creature-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "creature-1"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestCreatureGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	mockService.EXPECT().
		GetByID("creature-1").
		Return(nil, nil)

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/creature/creature-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "creature-1"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestCreatureGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	expected := &m.Creature{ID: "creature-1", Name: "Goblin"}

	mockService.EXPECT().
		GetByID("creature-1").
		Return(expected, nil)

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/creature/creature-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "creature-1"})
	rr := httptest.NewRecorder()

	controller.GetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body m.Creature
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.ID != "creature-1" {
		t.Fatalf("expected creature-1, got %s", body.ID)
	}
}

// --- Save ---

func TestCreatureSave_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/creature", bytes.NewBuffer([]byte("invalid json")))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreatureSave_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	creature := m.Creature{Name: "Goblin"}
	bodyBytes, _ := json.Marshal(creature)

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(nil, h.AssertError())

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/creature", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestCreatureSave_BadRequestChangingTimestamps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	creature := m.Creature{Name: "Goblin"}
	bodyBytes, _ := json.Marshal(creature)

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(nil, m.NewBadRequestChangingTimestampsError())

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/creature", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreatureSave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	now := time.Now()
	creature := m.Creature{Name: "Goblin"}
	bodyBytes, _ := json.Marshal(creature)

	expected := &m.Creature{ID: "generated-id", Name: "Goblin", CreatedAt: &now, UpdatedAt: &now}

	mockService.EXPECT().
		Save(gomock.Any()).
		Return(expected, nil)

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost, "/creature", bytes.NewBuffer(bodyBytes))
	rr := httptest.NewRecorder()

	controller.Save(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body m.Creature
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if body.ID != "generated-id" {
		t.Fatalf("expected generated-id, got %s", body.ID)
	}
}

// --- Delete ---

func TestCreatureDelete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	mockService.EXPECT().
		Delete("creature-1").
		Return(h.AssertError())

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodDelete, "/creature/creature-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "creature-1"})
	rr := httptest.NewRecorder()

	controller.Delete(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestCreatureDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	mockService.EXPECT().
		Delete("creature-1").
		Return(nil)

	controller := NewCreatureController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodDelete, "/creature/creature-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "creature-1"})
	rr := httptest.NewRecorder()

	controller.Delete(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}

// --- HandleSubrouter ---

func TestHandleCreatureSubrouter_WiresRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockCreatureService(ctrl)

	controller := NewCreatureController(mockLogger, mockService)
	router := mux.NewRouter()
	controller.HandleSubrouter(router)

	// GET /creature
	mockService.EXPECT().GetAll().Return(&[]m.Creature{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/creature", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /creature expected %d, got %d", http.StatusOK, rr.Code)
	}

	// GET /creature/{id}
	mockService.EXPECT().GetByID("creature-1").Return(&m.Creature{ID: "creature-1", Name: "Goblin"}, nil)
	req = httptest.NewRequest(http.MethodGet, "/creature/creature-1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /creature/creature-1 expected %d, got %d", http.StatusOK, rr.Code)
	}

	// POST /creature
	mockService.EXPECT().Save(gomock.Any()).Return(&m.Creature{ID: "creature-1", Name: "Goblin"}, nil)
	bodyBytes, _ := json.Marshal(m.Creature{Name: "Goblin"})
	req = httptest.NewRequest(http.MethodPost, "/creature", bytes.NewBuffer(bodyBytes))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /creature expected %d, got %d", http.StatusOK, rr.Code)
	}

	// DELETE /creature/{id}
	mockService.EXPECT().Delete("creature-1").Return(nil)
	req = httptest.NewRequest(http.MethodDelete, "/creature/creature-1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("DELETE /creature/creature-1 expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}
