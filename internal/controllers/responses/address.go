package responses

import (
	"project-skbackend/internal/models/helper"

	"github.com/google/uuid"
)

type (
	Address struct {
		helper.Model
		UserID      uuid.UUID `json:"-"`
		Name        string    `json:"name"`
		Address     string    `json:"address"`
		Description string    `json:"description"`
		Note        string    `json:"note"`
		Landmark    string    `json:"landmark"`
		Longitude   float64   `json:"langitude"`
		Latitude    float64   `json:"latitude"`
	}
)
