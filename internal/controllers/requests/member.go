package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utmath"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateMember struct {
		User           CreateUser        `json:"user"`
		Caregiver      *CreateCaregiver  `json:"caregiver"`
		Height         float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight         float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		FirstName      string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName       string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		Gender         consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
		OrganizationID *uuid.UUID        `json:"organization_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		IllnessID      []uuid.UUID       `json:"illness_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID      []uuid.UUID       `json:"allergy_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
	}

	UpdateMember struct {
		User           UpdateUser        `json:"user"`
		Caregiver      *UpdateCaregiver  `json:"caregiver"`
		Height         float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight         float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		FirstName      string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName       string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		Gender         consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
		OrganizationID *uuid.UUID        `json:"organization_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		IllnessID      []uuid.UUID       `json:"illness_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID      []uuid.UUID       `json:"allergy_id" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
	}
)

func (req *CreateMember) ToModel(
	user models.User,
	caregiver *models.Caregiver,
	allergies []*models.MemberAllergy,
	illnesses []*models.MemberIllness,
	organization *models.Organization,
) *models.Member {
	member := models.Member{
		User:         user,
		Caregiver:    caregiver,
		Allergy:      allergies,
		Illness:      illnesses,
		Organization: organization,
		BMI:          utmath.BMICalculation(req.Weight, req.Height),
	}

	if err := copier.Copy(&member, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &member
}

func (req *UpdateMember) ToModel(
	member models.Member,
	user models.User,
	caregiver *models.Caregiver,
	allergies []*models.MemberAllergy,
	illnesses []*models.MemberIllness,
	organization *models.Organization,
) *models.Member {
	if err := copier.CopyWithOption(&member, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	member.User = user
	member.Caregiver = caregiver
	member.Allergy = allergies
	member.Illness = illnesses
	member.Organization = organization
	member.BMI = utmath.BMICalculation(req.Weight, req.Height)

	return &member
}
