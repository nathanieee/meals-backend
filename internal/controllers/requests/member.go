package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateMemberRequest struct {
		User           CreateUserRequest       `json:"user"`
		Caregiver      *CreateCaregiverRequest `json:"caregiver"`
		Height         float64                 `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight         float64                 `json:"weight" gorm:"not null" binding:"required" example:"150"`
		FirstName      string                  `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName       string                  `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		Gender         consttypes.Gender       `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth    time.Time               `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
		OrganizationID *uuid.UUID              `json:"organization_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		IllnessID      []uuid.UUID             `json:"illness_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID      []uuid.UUID             `json:"allergy_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
	}

	UpdateMemberRequest struct {
		Height      float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight      float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		FirstName   string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth time.Time         `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
		IllnessID   []uuid.UUID       `json:"illness_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID   []uuid.UUID       `json:"allergy_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
	}
)

func (req *CreateMemberRequest) ToModel(
	user models.User,
	caregiver models.Caregiver,
	allergies []*models.MemberAllergy,
	illness []*models.MemberIllness,
	organization *models.Organization,
) *models.Member {
	member := models.Member{
		User:         user,
		Caregiver:    &caregiver,
		Allergy:      allergies,
		Illness:      illness,
		Organization: organization,
	}

	if err := copier.Copy(&member, &req); err != nil {
		logger.LogError(err)
		return nil
	}

	return &member
}
