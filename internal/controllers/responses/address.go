package responses

import (
	"project-skbackend/internal/models/helper"
)

type (
	Address struct {
		helper.Model
		Name      string  `json:"name"`
		Address   string  `json:"address"`
		Note      string  `json:"note"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}
)
