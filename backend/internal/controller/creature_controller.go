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

type CreatureController struct {
	logger          l.Logger
	creatureService s.CreatureService
}

func NewCreatureController(
	logger l.Logger,
	creatureService s.CreatureService,
) *CreatureController {
	controllerLogger := logger.WithFields("file_name", "creature_controller.go", "class", "CreatureController")

	return &CreatureController{
		logger:          controllerLogger,
		creatureService: creatureService,
	}
}

func (c *CreatureController) HandleSubrouter(r *mux.Router) {
	c.logger.Debug("Handle Subrouter.")

	router := r.PathPrefix("/creature").Subrouter()

	c.logger.Trace("Handle GET /creature endpoint.")
	router.HandleFunc("", c.GetAll).Methods("GET")

	c.logger.Trace("Handle GET /creature/{id} endpoint.")
	router.HandleFunc("/{id}", c.GetByID).Methods("GET")

	c.logger.Trace("Handle POST /creature endpoint.")
	router.HandleFunc("", c.Save).Methods("POST")

	c.logger.Trace("Handle DELETE /creature/{id} endpoint.")
	router.HandleFunc("/{id}", c.Delete).Methods("DELETE")
}

func (c *CreatureController) GetAll(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get All Creatures endpoint hit.")

	results, err := c.creatureService.GetAll()
	if err != nil {
		c.logger.Error("Failed to get all creatures.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		c.logger.Error("Failed to encode creatures.", "error", err.Error())
	}

	c.logger.Debug("Completed Get All Creatures endpoint.")
}

func (c *CreatureController) GetByID(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Get Creature By ID endpoint hit.")

	vars := mux.Vars(r)
	id := vars["id"]

	c.logger.Trace("Get creature by id.", "id", id)
	result, err := c.creatureService.GetByID(id)
	if err != nil {
		c.logger.Error("Failed to get creature by id.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if result == nil {
		c.logger.Debug("Creature not found.", "id", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		c.logger.Error("Failed to encode creature.", "error", err.Error())
	}

	c.logger.Debug("Completed Get Creature By ID endpoint.")
}

func (c *CreatureController) Save(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Save Creature endpoint hit.")

	var creature m.Creature

	c.logger.Trace("Decode request body.")
	if err := json.NewDecoder(r.Body).Decode(&creature); err != nil {
		c.logger.Warn("Failed to decode request body, bad request from user.")
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	c.logger.Trace("Call service to save creature.")
	result, err := c.creatureService.Save(creature)
	if err != nil {
		c.logger.Error("Failed to save creature.", "error", err.Error())
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
		c.logger.Error("Failed to encode saved creature.", "error", err.Error())
	}

	c.logger.Debug("Completed Save Creature endpoint.")
}

func (c *CreatureController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug("Delete Creature endpoint hit.")

	vars := mux.Vars(r)
	id := vars["id"]

	c.logger.Trace("Call service to delete creature.", "id", id)
	if err := c.creatureService.Delete(id); err != nil {
		c.logger.Error("Failed to delete creature.", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	c.logger.Debug("Completed Delete Creature endpoint.")
}
