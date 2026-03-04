package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DieType int

const (
	D2  DieType = 2
	D4  DieType = 4
	D6  DieType = 6
	D8  DieType = 8
	D10 DieType = 10
	D12 DieType = 12
	D20 DieType = 20
)

// DiceExpression represents a dice roll with an optional flat modifier.
// Examples: 2d6 = {NumDice:2, DieType:D6, Modifier:0}
//
//	1d6+3 = {NumDice:1, DieType:D6, Modifier:3}
//	flat 4 = {NumDice:0, DieType:0, Modifier:4}
type DiceExpression struct {
	NumDice  int     `json:"num_dice"`
	DieType  DieType `json:"die_type"`
	Modifier int     `json:"modifier"`
}

// Value implements driver.Valuer so DiceExpression can be written to a jsonb column.
func (d DiceExpression) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DiceExpression: %w", err)
	}
	return string(b), nil
}

// Scan implements sql.Scanner so DiceExpression can be read from a jsonb column.
func (d *DiceExpression) Scan(src any) error {
	var raw []byte
	switch v := src.(type) {
	case string:
		raw = []byte(v)
	case []byte:
		raw = v
	default:
		return fmt.Errorf("unsupported type for DiceExpression: %T", src)
	}

	if err := json.Unmarshal(raw, d); err != nil {
		return fmt.Errorf("failed to unmarshal DiceExpression: %w", err)
	}
	return nil
}
