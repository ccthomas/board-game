package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/ccthomas/board-game/internal/helper"
	s "github.com/ccthomas/board-game/internal/service/mock"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func TestGetHealth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	mockService.
		EXPECT().
		GetDatabaseVersion().
		Return("1.0.0")

	controller := NewHealthController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	controller.GetHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if body["service"] != "Healthy" {
		t.Fatalf("expected service Healthy, got %v", body["service"])
	}

	if body["database"] != "1.0.0" {
		t.Fatalf("expected database 1.0.0, got %v", body["database"])
	}

	if body["timestamp"] == nil {
		t.Fatalf("expected timestamp to exist")
	}
}

func TestGetHealth_DatabaseErrorString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	mockService.
		EXPECT().
		GetDatabaseVersion().
		Return("database unreachable")

	controller := NewHealthController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	controller.GetHealth(rr, req)

	var body map[string]interface{}
	_ = json.Unmarshal(rr.Body.Bytes(), &body)

	if body["database"] != "database unreachable" {
		t.Fatalf("expected error string, got %v", body["database"])
	}
}

func TestHandleHealthSubrouter_WiresRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	mockService.
		EXPECT().
		GetDatabaseVersion().
		Return("1.0.0")

	controller := NewHealthController(mockLogger, mockService)

	router := mux.NewRouter()
	controller.HandleSubrouter(router)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}
}
