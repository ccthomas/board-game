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

type DamageTypeController struct {
	logger            l.Logger
	damageTypeService s.DamageTypeService
}

func NewDamageTypeController(
	logger l.Logger,
	damageTypeService s.DamageTypeService,
) *DamageTypeController {
	controllerLogger := logger.WithFields("file_name", "damage_type_controller.go", "class", "DamageTypeController")

	return &DamageTypeController{
		logger:            controllerLogger,
		damageTypeService: damageTypeService,
	}
}

func (c *DamageTypeController) HandleSubrouter(r *mux.Router) {
	c.logger.Debug("Handle Subrouter.")

	router := r.PathPrefix("/damage-type").Subrouter()

	c.logger.Trace("Handle GET /damage-type endpoint.")
	router.HandleFunc("", c.GetAll).Methods("GET")

	c.logger.Trace("Handle GET /damage-type/{id} endpoint.")
	router.HandleFunc("/{id}", c.GetByID).Methods("GET")

	c.logger.Trace("Handle POST /damage-type endpoint.")
	router.HandleFunc("", c.Save).Methods("POST")

	c.logger.Trace("Handle DELETE /damage-type/{id} endpoint.")
	router.HandleFunc("/{id}", c.Delete).Methods("DELETE")
}

func (c *DamageTypeController) GetAll(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get All Damage Types endpoint hit.")

	results, err := c.damageTypeService.GetAll()
	if err != nil {
		c.logger.Error("Failed to get all damage types.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		c.logger.Error("Failed to encode damage types.", "error", err.Error())
	}

	c.logger.Debug("Completed Get All Damage Types endpoint.")
}

func (c *DamageTypeController) GetByID(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get Damage Type By ID endpoint hit.")

	vars := mux.Vars(r)
	id := vars["id"]

	c.logger.Trace("Get damage type by id.", "id", id)
	result, err := c.damageTypeService.GetByID(id)
	if err != nil {
		c.logger.Error("Failed to get damage type by id.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if result == nil {
		c.logger.Debug("Damage type not found.", "id", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		c.logger.Error("Failed to encode damage type.", "error", err.Error())
	}

	c.logger.Debug("Completed Get Damage Type By ID endpoint.")
}

func (c *DamageTypeController) Save(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Save Damage Type endpoint hit.")

	var damageType m.DamageType

	c.logger.Trace("Decode request body.")
	if err := json.NewDecoder(r.Body).Decode(&damageType); err != nil {
		c.logger.Warn("Failed to decode request body, bad request from user.")
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	c.logger.Trace("Call service to save damage type.")
	result, err := c.damageTypeService.Save(damageType)
	if err != nil {
		c.logger.Error("Failed to save damage type.", "error", err.Error())
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
		c.logger.Error("Failed to encode saved damage type.", "error", err.Error())
	}

	c.logger.Debug("Completed Save Damage Type endpoint.")
}

func (c *DamageTypeController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Delete Damage Type endpoint hit.")

	vars := mux.Vars(r)
	id := vars["id"]

	c.logger.Trace("Call service to delete damage type.", "id", id)
	if err := c.damageTypeService.Delete(id); err != nil {
		c.logger.Error("Failed to delete damage type.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	c.logger.Debug("Completed Delete Damage Type endpoint.")
}
