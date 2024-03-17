package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
)

type (
	Caregiver struct {
		base.Model

		User User `json:"user"`

		Gender      consttypes.Gender `json:"gender"`
		FirstName   string            `json:"first_name"`
		LastName    string            `json:"last_name"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth"`
	}
)
