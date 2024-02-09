package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Organization struct {
		helper.Model
		User User                        `json:"user"`
		Type consttypes.OrganizationType `json:"type"`
		Name string                      `json:"name"`
	}
)
