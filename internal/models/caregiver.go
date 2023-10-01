package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	Caregiver struct {
		helper.Model
		UserID      uuid.UUID         `json:"userID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User        User              `json:"user"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		FirstName   string            `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		DateOfBirth time.Time         `json:"date" gorm:"not null" binding:"required" example:"2000-10-20"`
	}
)
