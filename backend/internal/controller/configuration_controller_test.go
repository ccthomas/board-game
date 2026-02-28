package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/ccthomas/board-game/internal/helper"
	m "github.com/ccthomas/board-game/internal/model"
	s "github.com/ccthomas/board-game/internal/service/mock"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func TestRunDatabaseMigration_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockConfigurationService(ctrl)

	controller := NewConfigurationController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost,
		"/configuration/database/migrations",
		bytes.NewBuffer([]byte("invalid json")),
	)

	rr := httptest.NewRecorder()

	controller.RunDatabaseMigration(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestRunDatabaseMigration_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockConfigurationService(ctrl)

	config := m.MigrationConfiguration{
		Command: m.MigrationUp,
	}

	bodyBytes, _ := json.Marshal(config)

	mockService.
		EXPECT().
		RunDatabaseMigration(gomock.Any()).
		Return(h.AssertError())

	controller := NewConfigurationController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost,
		"/configuration/database/migrations",
		bytes.NewBuffer(bodyBytes),
	)

	rr := httptest.NewRecorder()

	controller.RunDatabaseMigration(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestRunDatabaseMigration_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockConfigurationService(ctrl)

	config := m.MigrationConfiguration{
		Command: m.MigrationUp,
	}

	bodyBytes, _ := json.Marshal(config)

	mockService.
		EXPECT().
		RunDatabaseMigration(gomock.Any()).
		Return(nil)

	controller := NewConfigurationController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodPost,
		"/configuration/database/migrations",
		bytes.NewBuffer(bodyBytes),
	)

	rr := httptest.NewRecorder()

	controller.RunDatabaseMigration(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestHandleConfigurationSubrouter_WiresRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockConfigurationService(ctrl)

	mockService.
		EXPECT().
		RunDatabaseMigration(gomock.Any()).
		Return(nil)

	controller := NewConfigurationController(mockLogger, mockService)

	router := mux.NewRouter()
	controller.HandleSubrouter(router)

	config := m.MigrationConfiguration{
		Command: m.MigrationUp,
	}

	bodyBytes, _ := json.Marshal(config)

	req := httptest.NewRequest(http.MethodPost,
		"/configuration/database/migrations",
		bytes.NewBuffer(bodyBytes),
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}
}
