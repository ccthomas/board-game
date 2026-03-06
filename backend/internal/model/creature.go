package model

import "time"

// Creature is the static, reusable definition of a creature.
// Shared across games — never mutated at runtime.
type Creature struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	HealthPoints int            `json:"health_points"`
	Defence      DiceExpression `json:"defence"`
	Initiative   int            `json:"initiative"`
	Movement     int            `json:"movement"`
	ActionCount  int            `json:"action_count"` // number of actions this creature may take per turn
	CreatedAt    *time.Time     `json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at"`

	// transient
	Abilities []AbilitySlot `json:"abilities"` // ordered ascending by RollThreshold
}

type AbilitySlot struct {
	// transient
	Ability Ability `json:"ability"`
	// database column
	AbilityID  string `json:"ability_id"`
	CreatureID string `json:"creature_id"`

	RollThreshold int        `json:"roll_threshold"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
