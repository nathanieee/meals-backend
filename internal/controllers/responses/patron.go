package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Patron struct {
		base.Model

		User User `json:"user"`

		Type consttypes.PatronType `json:"type"`
		Name string                `json:"name"`

		Donations []Donation `json:"donations"`
	}
)
