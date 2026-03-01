// Package controller...TODO
package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	l "github.com/ccthomas/board-game/internal/logger"
	m "github.com/ccthomas/board-game/internal/model"
	s "github.com/ccthomas/board-game/internal/service"
	"github.com/gorilla/mux"
)

type ConfigurationController struct {
	logger               l.Logger
	configurationService s.ConfigurationService
}

func NewConfigurationController(
	logger l.Logger,
	configurationService s.ConfigurationService,
) *ConfigurationController {
	controllerLogger := logger.WithFields("file_name", "configuration_controller.go", "class", "ConfigurationController")

	return &ConfigurationController{
		logger:               controllerLogger,
		configurationService: configurationService,
	}
}

func (c *ConfigurationController) HandleSubrouter(r *mux.Router) {
	c.logger.Debug("Handle Subrounder.")

	router := r.PathPrefix("/configuration").Subrouter()

	c.logger.Trace("Handle POST /database/migration endpoint.")
	router.HandleFunc("/database/migrations", c.RunDatabaseMigration).Methods("POST")
}

func (c *ConfigurationController) RunDatabaseMigration(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Run Database Migration endpoint hit.")

	var config m.MigrationConfiguration

	c.logger.Trace("Decode request body.")
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		c.logger.Warn("Failed to decode request body, bad request from user.")

		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	c.logger.Trace("Call service to run database migration.")
	err := c.configurationService.RunDatabaseMigration(&config)
	if err != nil {
		c.logger.Error("Failed to run database migration:", "error", err.Error())

		if badReqErr, ok := errors.AsType[*m.BadMigrationCommandRequestError](err); ok {
			http.Error(w,
				http.StatusText(http.StatusUnprocessableEntity),
				http.StatusUnprocessableEntity,
			)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(badReqErr.Error()); err != nil {
				c.logger.Error("Failed to encode health data.", "error", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

		} else {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	c.logger.Debug("Completed Run Database Migration endpoint.")
}
