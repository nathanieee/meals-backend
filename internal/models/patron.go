package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Patron struct {
		helper.Model
		UserID uuid.UUID             `json:"user_id" gorm:"not null" binding:"required"`
		User   User                  `json:"user"`
		Type   consttypes.PatronType `json:"type" gorm:"not null; type:patron_type_enum" binding:"required"`
		Name   string                `json:"name" gorm:"not null" binding:"required" example:"Anonymus"`
	}
)
