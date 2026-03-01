// Package repository...TODO
package repository

import (
	"github.com/ccthomas/board-game/internal/model"
)

type DamageTypeRepository interface {
	GetAll() (*[]model.DamageType, error)
	GetByID(id string) (*model.DamageType, error)
	Upsert(d model.DamageType) error
}
