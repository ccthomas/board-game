package model

import "time"

type EffectType string

const (
	Damage  EffectType = "damage"
	Defence EffectType = "defence"
	Healing EffectType = "healing"
)

// AbilityEffect describes a single effect that resolves when an ability is used.
// An ability may have multiple effects — e.g. a lifedrain has one damage effect
// targeting an enemy and one healing effect targeting self.
type AbilityEffect struct {
	ID         string          `json:"id"`
	AbilityID  string          `json:"ability_id"`
	Expression DiceExpression  `json:"expression"` // Store as jsonb
	EffectType EffectType      `json:"effect_type"`
	Alignment  TargetAlignment `json:"alignment"` // who this specific effect targets
	CreatedAt  *time.Time      `json:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at"`

	// transient
	DamageType *DamageType `json:"damage_type"` // only relevant if Damage is set
	// database column
	DamageTypeID string
}

type Ability struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Pattern   TargetingPattern `json:"pattern"` // how targets are selected
	Range     int              `json:"range"`   // in tiles; 0 = self/adjacent only
	CreatedAt *time.Time       `json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at"`
	DeletedAt *time.Time       `json:"deleted_at"`

	// transient
	Effects []AbilityEffect `json:"effects"` // resolved in order; multiple effects support lifedrain etc.
}
