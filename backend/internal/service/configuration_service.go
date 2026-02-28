// Package service...TODO
package service

import (
	"errors"
	"fmt"

	d "github.com/ccthomas/board-game/internal/database"
	l "github.com/ccthomas/board-game/internal/logger"
	m "github.com/ccthomas/board-game/internal/model"
)

type ConfigurationService interface {
	RunDatabaseMigration(config *m.MigrationConfiguration) error
}

type ConfigurationServiceImpl struct {
	logger   l.Logger
	database d.Database
}

func NewConfigurationServiceImpl(logger l.Logger, database d.Database) *ConfigurationServiceImpl {
	serviceLogger := logger.WithFields(
		"file_name", "configuration_service.go",
		"class_name", "ConfigurationService",
	)

	return &ConfigurationServiceImpl{
		logger:   serviceLogger,
		database: database,
	}
}

func (s *ConfigurationServiceImpl) RunDatabaseMigration(config *m.MigrationConfiguration) error {
	s.logger.Debug("RunDatabaseMigration invoked.")

	s.logger.Trace("Attempting to get database connection.")
	db, err := s.database.GetConnection()
	if err != nil {
		s.logger.Error("Failed to get database connection.", "error", err.Error())
		return err
	}

	s.logger.Trace("Evaluating migration command.",
		"command", config.Command,
		"quantity", config.Quantity,
	)

	switch config.Command {

	case m.MigrationDown:
		s.logger.Debug("Executing MigrationDown command.")

		if config.Quantity == nil {
			s.logger.Trace("Running full MigrationUp (no quantity provided).")

			if err := s.database.MigrationUp(db); err != nil {
				s.logger.Error("MigrationUp failed.", "error", err.Error())
				return err
			}

		} else {
			steps := int8(*config.Quantity) * -1

			s.logger.Trace("Running MigrationSteps.",
				"steps", steps,
			)

			if err := s.database.MigrationSteps(db, steps); err != nil {
				s.logger.Error("MigrationSteps failed.", "error", err.Error())
				return err
			}
		}

	case m.MigrationUp:
		s.logger.Debug("Executing MigrationUp command.")

		if config.Quantity == nil {
			s.logger.Trace("Running full MigrationDown (no quantity provided).")

			if err := s.database.MigrationDown(db); err != nil {
				s.logger.Error("MigrationDown failed.", "error", err.Error())
				return err
			}

		} else {
			steps := int8(*config.Quantity)

			s.logger.Trace("Running MigrationSteps.",
				"steps", steps,
			)

			if err := s.database.MigrationSteps(db, steps); err != nil {
				s.logger.Error("MigrationSteps failed.", "error", err.Error())
				return err
			}
		}

	default:
		msg := fmt.Sprintf("database migration command unknown: %s", config.Command)
		s.logger.Warn("Unknown migration command received.",
			"command", config.Command,
		)
		return errors.New(msg)
	}

	s.logger.Debug("RunDatabaseMigration completed successfully.")
	return nil
}
