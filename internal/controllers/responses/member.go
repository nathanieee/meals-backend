package responses

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	MemberResponse struct {
		ID           uuid.UUID              `json:"id"`
		Email        string                 `json:"email" example:"email@email.com"`
		Role         consttypes.UserRole    `json:"role" gorm:"not null" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Caregiver    *CaregiverResponse     `json:"caregiver,omitempty"`
		Organization *OrganizationResponse  `json:"organization,omitempty"`
		Illness      *MemberIllnessResponse `json:"illness,omitempty"`
		Allergy      *MemberAllergyResponse `json:"allergy,omitempty"`
		Height       float64                `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight       float64                `json:"weight" gorm:"not null" binding:"required" example:"150"`
		BMI          float64                `json:"BMI" gorm:"not null" binding:"required" example:"19"`
		FirstName    string                 `json:"firstName" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName     string                 `json:"lastName" gorm:"not null" binding:"required" example:"Vince"`
		Gender       consttypes.Gender      `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth  time.Time              `json:"date" gorm:"not null" binding:"required" example:"2000-10-20"`
	}

	MemberIllnessResponse struct {
		ID        uuid.UUID      `json:"id"`
		IllnessID uuid.UUID      `json:"illnessID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Illness   models.Illness `json:"illness"`
	}

	MemberAllergyResponse struct {
		ID        uuid.UUID      `json:"id"`
		AllergyID uuid.UUID      `json:"allergyID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Allergy   models.Allergy `json:"allergy"`
	}
)
