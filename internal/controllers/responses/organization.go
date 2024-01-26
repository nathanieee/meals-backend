package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Organization struct {
		helper.Model
		UserID uuid.UUID                   `json:"-"`
		User   User                        `json:"user"`
		Type   consttypes.OrganizationType `json:"type"`
		Name   string                      `json:"name"`
	}
)
