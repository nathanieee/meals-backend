package requests

import (
	"project-skbackend/packages/consttypes"
	"time"
)

type (
	CreateCaregiverRequest struct {
		CreateUserRequest
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		FirstName   string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		DateOfBirth time.Time         `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
	}
)
