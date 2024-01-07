package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/custom"

	"github.com/google/uuid"
)

type (
	CaregiverResponse struct {
		helper.Model
		UserID      uuid.UUID         `json:"-"`
		User        UserResponse      `json:"user"`
		Gender      consttypes.Gender `json:"gender"`
		FirstName   string            `json:"first_name"`
		LastName    string            `json:"last_name"`
		DateOfBirth custom.CDT_DATE   `json:"date_of_birth"`
	}
)
