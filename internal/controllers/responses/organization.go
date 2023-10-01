package responses

import (
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	OrganizationResponse struct {
		ID   uuid.UUID                   `json:"id"`
		Type consttypes.OrganizationType `json:"type" gorm:"not null" binding:"required" example:"Orphanage"`
		Name string                      `json:"name" gorm:"not null" binding:"required" example:"Panti Jompo Syailendra"`
	}
)
