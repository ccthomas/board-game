// Package service...TODO
package service

import (
	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
)

type HealthService interface {
	GetDatabaseVersion() string
}

type HealthServiceImpl struct {
	logger   l.Logger
	database d.Database
}

func NewHealthServiceImpl(
	logger l.Logger,
	database d.Database,
) *HealthServiceImpl {
	serviceLogger := logger.WithFields(
		"file_name", "health_service.go",
		"class_name", "HealthService",
	)

	return &HealthServiceImpl{
		logger:   serviceLogger,
		database: database,
	}
}

func (s *HealthServiceImpl) GetDatabaseVersion() string {
	s.logger.Debug("GetDatabaseVersion invoked.")

	s.logger.Trace("Calling database.Version().")
	version, err := s.database.Version()
	if err != nil {
		s.logger.Error("Failed to retrieve database version.",
			"error", err.Error(),
		)
		return err.Error()
	}

	if version == nil {
		s.logger.Warn("Database version returned nil.")
		return "unknown"
	}

	s.logger.Debug("Successfully retrieved database version.",
		"version", *version,
	)

	return *version
}
