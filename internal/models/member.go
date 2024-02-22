package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Member struct {
		helper.Model
		UserID         uuid.UUID         `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User           User              `json:"user"`
		CaregiverID    *uuid.UUID        `json:"caregiver_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		Caregiver      *Caregiver        `json:"caregiver,omitempty"`
		OrganizationID *uuid.UUID        `json:"organization_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		Organization   *Organization     `json:"organization,omitempty"`
		Illnesses      []*MemberIllness  `json:"illnesses,omitempty"`
		Allergies      []*MemberAllergy  `json:"allergies,omitempty"`
		Height         float64           `json:"height" gorm:"required" example:"100"`
		Weight         float64           `json:"weight" gorm:"required" example:"150"`
		BMI            float64           `json:"bmi" gorm:"required;type:decimal(10,2)" example:"19"`
		FirstName      string            `json:"first_name" gorm:"required" example:"Jonathan"`
		LastName       string            `json:"last_name" gorm:"required" example:"Vince"`
		Gender         consttypes.Gender `json:"gender" gorm:"required; type:gender_enum" example:"Male"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth" gorm:"required" example:"2000-10-20"`
	}

	MemberIllness struct {
		helper.Model
		MemberID  uuid.UUID `json:"member_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		IllnessID uuid.UUID `json:"illness_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Illness   Illness   `json:"illness"`
	}

	MemberAllergy struct {
		helper.Model
		MemberID  uuid.UUID `json:"member_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID uuid.UUID `json:"allergy_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Allergy   Allergy   `json:"allergy"`
	}
)

func (m *Member) ToResponse() *responses.Member {
	mres := responses.Member{}

	if err := copier.CopyWithOption(&mres, &m, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &mres
}
