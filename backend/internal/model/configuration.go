package model

type MigrationCommand string

const (
	MigrationUp   MigrationCommand = "up"
	MigrationDown MigrationCommand = "down"
)

type MigrationConfiguration struct {
	Command  MigrationCommand `json:"command"`
	Quantity *int8            `json:"quantity"`
}
