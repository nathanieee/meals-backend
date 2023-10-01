package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	Member struct {
		helper.Model
		UserID         uuid.UUID         `json:"userID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User           User              `json:"user"`
		CaregiverID    uuid.UUID         `json:"caregiverID" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Caregiver      *Caregiver        `json:"caregiver,omitempty"`
		OrganizationID uuid.UUID         `json:"organizationID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Organization   *Organization     `json:"organization,omitempty"`
		Illness        *MemberIllness    `json:"illness,omitempty"`
		Allergy        *MemberAllergy    `json:"allergy,omitempty"`
		Height         float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight         float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		BMI            float64           `json:"BMI" gorm:"not null" binding:"required" example:"19"`
		FirstName      string            `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName       string            `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		Gender         consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth    time.Time         `json:"date" gorm:"not null" binding:"required" example:"2000-10-20"`
	}

	MemberIllness struct {
		helper.Model
		MemberID  uuid.UUID `json:"memberID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		IllnessID uuid.UUID `json:"illnessID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Illness   Illness   `json:"illness"`
	}

	MemberAllergy struct {
		helper.Model
		MemberID  uuid.UUID `json:"memberID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID uuid.UUID `json:"allergyID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Allergy   Allergy   `json:"allergy"`
	}
)
