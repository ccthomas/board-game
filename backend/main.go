package main

import (
	"log"
	"net/http"
	"os"

	c "github.com/ccthomas/board-game/internal/controller"
	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/repository"
	s "github.com/ccthomas/board-game/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func loadEnvironmentVars() {
	// Load the .env file
	// If args are passed, they will be an array of .env files to load
	// TODO: Come up with better pattern for local development
	// TODO: Do Not deploy this logic to a production environment, without additonal safegaurds
	if len(os.Args) > 1 {
		err := godotenv.Overload(os.Args[1:]...)
		if err != nil {
			return
		}
	}

	// fmt.Println(os.Getenv("DB_CONTAINER_NAME"))
}

func main() {
	loadEnvironmentVars()

	// -------------------------
	// Database
	// -------------------------
	logger, err := l.NewLoggerSlog()
	if err != nil {
		log.Panicf("Could not conifgure logger: %s\b", err.Error())
	}

	mainLogger := logger.WithFields("file_name", "main.go")

	mainLogger.Info("==============================")
	mainLogger.Info("Board Game Backend starting...")

	// -------------------------
	// Database
	// -------------------------
	mainLogger.Debug("Configuring databases...")
	database := d.NewDatabasePostgres(logger)

	// -------------------------
	// Repository
	// -------------------------
	mainLogger.Debug("Configuring databases...")
	abilityRepo := repository.NewAbilityRepositoryPostgres(logger, database)
	creatureRepo := repository.NewCreatureRepositoryPostgres(logger, database)
	damageTypeRepo := repository.NewDamageTypeRepositoryPostgres(logger, database)

	// -------------------------
	// Services
	// -------------------------
	mainLogger.Debug("Configuring services...")
	configurationService := s.NewConfigurationServiceImpl(logger, database)
	damageTypeService := s.NewDamageTypeServiceImpl(logger, damageTypeRepo)
	healthService := s.NewHealthServiceImpl(logger, database)

	abilityService := s.NewAbilityServiceImpl(
		logger,
		abilityRepo,
		damageTypeService,
	)

	creatureService := s.NewCreatureServiceImpl(
		logger,
		creatureRepo,
		abilityService,
	)

	// -------------------------
	// Controllers
	// -------------------------
	mainLogger.Debug("Configuring controllers...")
	creatureController := c.NewCreatureController(logger, creatureService)
	configurationController := c.NewConfigurationController(logger, configurationService)
	abilityController := c.NewAbilityController(logger, abilityService)
	damageTypeController := c.NewDamageTypeController(logger, damageTypeService)
	healthController := c.NewHealthController(logger, healthService)

	// -------------------------
	// Router
	// -------------------------
	mainLogger.Debug("Configuring router...")
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	abilityController.HandleSubrouter(apiRouter)
	configurationController.HandleSubrouter(apiRouter)
	creatureController.HandleSubrouter(apiRouter)
	damageTypeController.HandleSubrouter(apiRouter)
	healthController.HandleSubrouter(apiRouter)

	// -------------------------
	// Server
	// -------------------------
	addr := ":80"

	handler := handleCors(mainLogger, router)

	mainLogger.Info("Listening & Serving service on port.", "post", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		mainLogger.Error("Failed to listen and serve api.")
		os.Exit(1)
	}
}

func handleCors(logger l.Logger, router *mux.Router) http.Handler {
	logger.Info("Configuring cors...")
	allowedOrigins := []string{
		"http://localhost:3000",
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	return cors(router)
}
