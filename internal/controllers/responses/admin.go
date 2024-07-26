package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs/ctdatatype"
)

type (
	Admin struct {
		base.Model

		User User `json:"user"`

		FirstName   string              `json:"first_name"`
		LastName    string              `json:"last_name"`
		Gender      consttypes.Gender   `json:"gender"`
		DateOfBirth ctdatatype.CDT_DATE `json:"date_of_birth"`
	}
)
