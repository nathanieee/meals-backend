package responses

import "project-skbackend/internal/models/base"

type (
	Order struct {
		base.Model

		Member Member `json:"member"`

		Meal Meal `json:"meal"`

		Status string `json:"status"`
	}
)
