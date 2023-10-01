package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	Admin struct {
		helper.Model
		UserID      uuid.UUID         `json:"userID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User        User              `json:"user"`
		FirstName   string            `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth time.Time         `json:"date" gorm:"not null" binding:"required"` // TODO - add an example on the dateofbirth based on what format is the date
	}
)
