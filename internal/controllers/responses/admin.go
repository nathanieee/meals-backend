package responses

import (
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	AdminResponse struct {
		ID                 uuid.UUID           `json:"id"`
		Email              string              `json:"email" example:"email@email.com"`
		Role               consttypes.UserRole `json:"role" gorm:"not null" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		FirstName          string              `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName           string              `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		Gender             consttypes.Gender   `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth        time.Time           `json:"date" gorm:"not null" binding:"required" example:"2000-10-20"`
		ConfirmationSentAt time.Time           `json:"confirmationSentAt"`
		ConfirmedAt        time.Time           `json:"confirmedAt"`
		CreatedAt          time.Time           `json:"createdAt,omitempty" example:"2023-01-01T15:01:00+00:00"`
		UpdatedAt          time.Time           `json:"updatedAt,omitempty" example:"2023-02-11T15:01:00+00:00"`
	}
)
