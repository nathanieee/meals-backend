package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Organization struct {
		base.Model

		User User `json:"user"`

		Type consttypes.OrganizationType `json:"type"`
		Name string                      `json:"name"`
	}
)
