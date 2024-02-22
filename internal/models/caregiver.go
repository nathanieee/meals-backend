package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"

	"github.com/google/uuid"
)

type (
	Caregiver struct {
		helper.Model
		UserID      uuid.UUID         `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User        User              `json:"user"`
		Gender      consttypes.Gender `json:"gender" gorm:"required; type:gender_enum" example:"Male"`
		FirstName   string            `json:"first_name" gorm:"required" example:"Jonathan"`
		LastName    string            `json:"last_name" gorm:"required" example:"Vince"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth" gorm:"required" example:"2000-12-30"`
	}
)
