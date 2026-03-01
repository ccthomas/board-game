package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/ccthomas/board-game/internal/helper"
	"github.com/ccthomas/board-game/internal/model"
	s "github.com/ccthomas/board-game/internal/service/mock"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func parseHealthBody(t *testing.T, rr *httptest.ResponseRecorder) (map[string]interface{}, map[string]interface{}) {
	t.Helper()

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	database, ok := body["database"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected database to be an object, got %v", body["database"])
	}

	return body, database
}

func TestGetHealth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	databaseHealth := model.DatabaseHealth{Status: "Healthy"}

	mockService.EXPECT().
		GetDatabaseHealth().
		Return(databaseHealth)

	controller := NewHealthController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	controller.GetHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	body, database := parseHealthBody(t, rr)

	if body["service"] != "Healthy" {
		t.Fatalf("expected service Healthy, got %v", body["service"])
	}

	if body["timestamp"] == nil {
		t.Fatal("expected timestamp to exist")
	}

	if database["health"] != databaseHealth.Status {
		t.Fatalf("expected database health %s, got %v", databaseHealth.Status, database["health"])
	}
}

func TestGetHealth_DatabaseUnhealthy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	databaseHealth := model.DatabaseHealth{Status: "database unreachable"}

	mockService.EXPECT().
		GetDatabaseHealth().
		Return(databaseHealth)

	controller := NewHealthController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	controller.GetHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	_, database := parseHealthBody(t, rr)

	if database["health"] != databaseHealth.Status {
		t.Fatalf("expected health %s, got %v", databaseHealth.Status, database["health"])
	}
}

func TestGetHealth_WithMigrationStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	pending := 2
	migrationStatus := &model.MigrationStatus{
		CurrentVersion: 1,
		LatestVersion:  3,
		Pending:        pending,
		Total:          3,
		Dirty:          false,
	}

	databaseHealth := model.DatabaseHealth{
		Status:          "Healthy",
		MigrationStatus: migrationStatus,
	}

	mockService.EXPECT().
		GetDatabaseHealth().
		Return(databaseHealth)

	controller := NewHealthController(mockLogger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	controller.GetHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	_, database := parseHealthBody(t, rr)

	migration, ok := database["migration_status"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected migration_status to be an object, got %v", database["migration_status"])
	}

	if int(migration["pending"].(float64)) != pending {
		t.Fatalf("expected pending %d, got %v", pending, migration["pending"])
	}
}

func TestHandleHealthSubrouter_WiresRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := h.NewDummyMockedLogger(ctrl)
	mockService := s.NewMockHealthService(ctrl)

	mockService.EXPECT().
		GetDatabaseHealth().
		Return(model.DatabaseHealth{Status: "Healthy"})

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
