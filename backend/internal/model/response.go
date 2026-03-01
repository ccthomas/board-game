package model

type HealthResponse struct {
	Service   string         `json:"service"`
	Database  DatabaseHealth `json:"database"`
	Timestamp interface{}    `json:"timestamp"`
}

type DatabaseHealth struct {
	Status          string           `json:"health"`
	Version         string           `json:"version"`
	MigrationStatus *MigrationStatus `json:"migration_status"`
}
