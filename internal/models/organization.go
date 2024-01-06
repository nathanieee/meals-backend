package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Organization struct {
		helper.Model
		UserID uuid.UUID                   `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User                        `json:"user"`
		Type   consttypes.OrganizationType `json:"type" gorm:"not null; type:organization_type_enum" binding:"required" example:"Orphanage"`
		Name   string                      `json:"name" gorm:"not null" binding:"required" example:"Panti Jompo Syailendra"`
	}
)
