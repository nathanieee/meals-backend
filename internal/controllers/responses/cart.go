package responses

import "project-skbackend/internal/models/helper"

type (
	Cart struct {
		helper.Model

		Meal      Meal       `json:"meal,omitempty"`
		Member    *Member    `json:"member,omitempty"`
		Caregiver *Caregiver `json:"caregiver,omitempty"`
		Quantity  uint       `json:"quantity,omitempty"`
	}
)
