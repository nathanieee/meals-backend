package models

import (
	"project-skbackend/packages/consttypes"

	"gorm.io/gorm"
)

type (
	Role struct {
		gorm.Model
		Name    consttypes.Role  `json:"name,omitempty"`
		LevelID consttypes.Level `json:"levelID,omitempty"`
		Level   Level            `json:"-"`
	}
)
