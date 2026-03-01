// Package service...TODO
package service

import (
	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
	"github.com/ccthomas/board-game/internal/model"
)

type HealthService interface {
	GetDatabaseHealth() model.DatabaseHealth
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

func (s *HealthServiceImpl) GetDatabaseHealth() model.DatabaseHealth {
	s.logger.Debug("GetDatabaseVersion invoked.")
	status := model.DatabaseHealth{
		Status:          "Healthy",
		Version:         "unknown",
		MigrationStatus: nil,
	}

	s.logger.Trace("Calling database.Version().")
	version, err := s.database.Version()
	if err != nil {
		s.logger.Error("Failed to retrieve database version.",
			"error", err.Error(),
		)
		status.Status = err.Error()
		return status
	}

	if version == nil {
		s.logger.Warn("Database version returned nil.")
		return status
	}

	status.Version = *version

	s.logger.Trace("Calling database.MigrationStatus().")
	migrationStatus, err := s.database.MigrationStatus()
	if err != nil {
		s.logger.Error("Failed to retrieve migration status.",
			"error", err.Error(),
		)
		status.Status = err.Error()
		return status
	}

	status.MigrationStatus = migrationStatus

	s.logger.Debug("Successfully retrieved database status.",
		"status", status,
	)

	return status
}
