// Package model...TODO
// TODO Create Error Response model to return to user with; message, timestamp, generated trace id
package model

import "fmt"

type BadMigrationCommandRequestError struct {
	message string
}

type BadRequestChangingTimestampsError struct {
	message string
}

func NewBadMigrationCommandRequestError(format string, a ...any) (err *BadMigrationCommandRequestError) {
	return &BadMigrationCommandRequestError{fmt.Sprintf(format, a...)}
}

func (m *BadMigrationCommandRequestError) Error() string {
	return m.message
}

func NewBadRequestChangingTimestampsError() (err *BadRequestChangingTimestampsError) {
	return &BadRequestChangingTimestampsError{"Bad Request: User cannot modify timestamp fields on data"}
}

func (m *BadRequestChangingTimestampsError) Error() string {
	return m.message
}
