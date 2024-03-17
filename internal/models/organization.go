package models

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Organization struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User      `json:"user"`

		Type consttypes.OrganizationType `json:"type" gorm:"required; type:organization_type_enum" example:"Orphanage"`
		Name string                      `json:"name" gorm:"required" example:"Panti Jompo Syailendra"`
	}
)
