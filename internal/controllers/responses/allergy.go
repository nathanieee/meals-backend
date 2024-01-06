package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	AllergyResponse struct {
		helper.Model
		Name        string               `json:"name"`
		Description string               `json:"description"`
		Allergens   consttypes.Allergens `json:"allergens"`
	}
)
