package model

import "time"

type DamageType struct {
	ID string `json:"id"`
	// TODO add business logic for unique name
	// As we prepare for an event driven world, database constrains should not be used
	// The database will be a reference to the truth (events), so validation should be in the BL Layer
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
