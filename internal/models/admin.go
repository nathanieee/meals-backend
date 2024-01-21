package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"

	"github.com/google/uuid"
)

type (
	Admin struct {
		helper.Model
		UserID      uuid.UUID         `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User        User              `json:"user"`
		FirstName   string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null; type:gender_enum" binding:"required" example:"Male"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-12-30"`
	}
)
