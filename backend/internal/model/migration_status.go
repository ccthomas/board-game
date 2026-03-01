package model

type MigrationStatus struct {
	CurrentVersion uint `json:"current_version"`
	LatestVersion  uint `json:"latest_version"`
	Pending        int  `json:"pending"`
	Total          int  `json:"total"`
	Dirty          bool `json:"dirty"`
}
