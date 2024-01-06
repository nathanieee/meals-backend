package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	AdminResponse struct {
		helper.Model
		UserID      uuid.UUID         `json:"-"`
		User        UserResponse      `json:"user"`
		FirstName   string            `json:"first_name"`
		LastName    string            `json:"last_name"`
		Gender      consttypes.Gender `json:"gender"`
		DateOfBirth time.Time         `json:"date_of_birth"`
	}
)
