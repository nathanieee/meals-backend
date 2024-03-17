package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Allergy struct {
		base.Model

		Name        string               `json:"name"`
		Description string               `json:"description"`
		Allergens   consttypes.Allergens `json:"allergens"`
	}
)
