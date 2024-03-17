package models

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Patron struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required"`
		User   User      `json:"user"`

		Type consttypes.PatronType `json:"type" gorm:"required; type:patron_type_enum"`
		Name string                `json:"name" gorm:"required" example:"Anonymus"`
	}
)
