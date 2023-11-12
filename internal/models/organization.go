package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Organization struct {
		helper.Model
		Type consttypes.OrganizationType `json:"type" gorm:"not null; type:organization_type_enum" binding:"required" example:"Orphanage"`
		Name string                      `json:"name" gorm:"not null" binding:"required" example:"Panti Jompo Syailendra"`
	}
)
