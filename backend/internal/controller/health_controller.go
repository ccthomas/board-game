// Package controller...TODO
package controller

import (
	"encoding/json"
	"net/http"
	"time"

	l "github.com/ccthomas/board-game/internal/logger"
	s "github.com/ccthomas/board-game/internal/service"
	"github.com/gorilla/mux"
)

type HealthController struct {
	logger        l.Logger
	healthService s.HealthService
}

func NewHealthController(logger l.Logger, healthService s.HealthService) *HealthController {
	controllerLogger := logger.WithFields("file_name", "health_controller.go", "class", "HealthController")

	return &HealthController{
		logger:        controllerLogger,
		healthService: healthService,
	}
}

func (c HealthController) HandleSubrouter(r *mux.Router) {
	c.logger.Debug("Handle Subrounter.")
	router := r.PathPrefix("/health").Subrouter()

	c.logger.Trace("Handle GET /health endpoint.")
	router.HandleFunc("", c.GetHealth).Methods("GET")
}

func (c *HealthController) GetHealth(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get Health endpoint hit.")
	databaseVersionOrError := c.healthService.GetDatabaseVersion()

	c.logger.Trace("Building health data.")
	currentTime := time.Now()
	data := map[string]interface{}{
		"service":   "Healthy",
		"database":  databaseVersionOrError,
		"timestamp": currentTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	c.logger.Trace("Encoding health data.", data)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		c.logger.Error("Failed to encode health data.", data)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	c.logger.Debug("Completed Get Health endpoint.")
}
