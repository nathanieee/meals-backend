package requests

import (
	"project-skbackend/packages/consttypes"
	"time"
)

type (
	CreateCaregiverRequest struct {
		CreateUserRequest
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		FirstName   string            `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		DateOfBirth time.Time         `json:"dateOfBirth" gorm:"not null" binding:"required" example:"2000-10-20"`
	}
)
