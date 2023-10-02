package requests

import (
	"project-skbackend/packages/consttypes"
	"time"
)

type (
	CreateMemberRequest struct {
		CreateUserRequest
		Height      float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight      float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		FirstName   string            `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth time.Time         `json:"date" gorm:"not null" binding:"required" example:"2000-10-20"`
	}

	UpdateMemberRequest struct {
		Height      float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight      float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		FirstName   string            `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth time.Time         `json:"date" gorm:"not null" binding:"required" example:"2000-10-20"`
	}
)
