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

type AbilityController struct {
	logger         l.Logger
	abilityService s.AbilityService
}

func NewAbilityController(
	logger l.Logger,
	abilityService s.AbilityService,
) *AbilityController {
	controllerLogger := logger.WithFields("file_name", "ability_controller.go", "class", "AbilityController")

	return &AbilityController{
		logger:         controllerLogger,
		abilityService: abilityService,
	}
}

func (c *AbilityController) HandleSubrouter(r *mux.Router) {
	c.logger.Debug("Handle Subrouter.")

	router := r.PathPrefix("/ability").Subrouter()

	c.logger.Trace("Handle GET /ability endpoint.")
	router.HandleFunc("", c.GetAll).Methods("GET")

	c.logger.Trace("Handle GET /ability/{id} endpoint.")
	router.HandleFunc("/{id}", c.GetByID).Methods("GET")

	c.logger.Trace("Handle POST /ability endpoint.")
	router.HandleFunc("", c.Save).Methods("POST")

	c.logger.Trace("Handle DELETE /ability/{id} endpoint.")
	router.HandleFunc("/{id}", c.Delete).Methods("DELETE")
}

func (c *AbilityController) GetAll(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get All Abilities endpoint hit.")

	results, err := c.abilityService.GetAll()
	if err != nil {
		c.logger.Error("Failed to get all abilities.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		c.logger.Error("Failed to encode abilities.", "error", err.Error())
	}

	c.logger.Debug("Completed Get All Abilities endpoint.")
}

func (c *AbilityController) GetByID(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get Ability By ID endpoint hit.")

	vars := mux.Vars(r)
	id := vars["id"]

	c.logger.Trace("Get ability by id.", "id", id)
	result, err := c.abilityService.GetByID(id)
	if err != nil {
		c.logger.Error("Failed to get ability by id.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if result == nil {
		c.logger.Debug("Ability not found.", "id", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		c.logger.Error("Failed to encode ability.", "error", err.Error())
	}

	c.logger.Debug("Completed Get Ability By ID endpoint.")
}

func (c *AbilityController) Save(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Save Ability endpoint hit.")

	var ability m.Ability

	c.logger.Trace("Decode request body.")
	if err := json.NewDecoder(r.Body).Decode(&ability); err != nil {
		c.logger.Warn("Failed to decode request body, bad request from user.")
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	c.logger.Trace("Call service to save ability.")
	result, err := c.abilityService.Save(ability)
	if err != nil {
		c.logger.Error("Failed to save ability.", "error", err.Error())
		if badReqErr, ok := errors.AsType[*m.BadRequestChangingTimestampsError](err); ok {
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(badReqErr.Error()); err != nil {
				c.logger.Error("Failed to encode error data.", "error", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		c.logger.Error("Failed to encode saved ability.", "error", err.Error())
	}

	c.logger.Debug("Completed Save Ability endpoint.")
}

func (c *AbilityController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Delete Ability endpoint hit.")

	vars := mux.Vars(r)
	id := vars["id"]

	c.logger.Trace("Call service to delete ability.", "id", id)
	if err := c.abilityService.Delete(id); err != nil {
		c.logger.Error("Failed to delete ability.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	c.logger.Debug("Completed Delete Ability endpoint.")
}
