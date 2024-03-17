package responses

import "project-skbackend/internal/models/base"

type (
	Cart struct {
		base.Model

		Meal Meal `json:"meal"`

		Member *Member `json:"member,omitempty"`

		Caregiver *Caregiver `json:"caregiver,omitempty"`

		Quantity uint `json:"quantity"`
	}
)
