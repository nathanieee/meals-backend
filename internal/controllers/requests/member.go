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
		User           CreateUser        `json:"user" form:"user" binding:"required"`
		Caregiver      *CreateCaregiver  `json:"caregiver" form:"caregiver" binding:"-"`
		Height         float64           `json:"height" form:"height" binding:"required"`
		Weight         float64           `json:"weight" form:"weight" binding:"required"`
		FirstName      string            `json:"first_name" form:"first_name" binding:"required"`
		LastName       string            `json:"last_name" form:"last_name" binding:"required"`
		Gender         consttypes.Gender `json:"gender" form:"gender" binding:"required"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth" form:"date_of_birth" binding:"required"`
		OrganizationID *uuid.UUID        `json:"organization_id,omitempty" form:"organization_id" binding:"-"`
		IllnessID      []*uuid.UUID      `json:"illness_id" form:"illness_id" binding:"-"`
		AllergyID      []*uuid.UUID      `json:"allergy_id" form:"allergy_id" binding:"-"`
	}

	UpdateMember struct {
		User           UpdateUser        `json:"user" form:"user" binding:"required"`
		Caregiver      *UpdateCaregiver  `json:"caregiver" form:"caregiver" binding:"-"`
		Height         float64           `json:"height" form:"height" binding:"required"`
		Weight         float64           `json:"weight" form:"weight" binding:"required"`
		FirstName      string            `json:"first_name" form:"first_name" binding:"required"`
		LastName       string            `json:"last_name" form:"last_name" binding:"required"`
		Gender         consttypes.Gender `json:"gender" form:"gender" binding:"required"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth" form:"date_of_birth" binding:"required"`
		OrganizationID *uuid.UUID        `json:"organization_id,omitempty" form:"organization_id" binding:"-"`
		IllnessID      []*uuid.UUID      `json:"illness_id" form:"illness_id" binding:"-"`
		AllergyID      []*uuid.UUID      `json:"allergy_id" form:"allergy_id" binding:"-"`
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
		Allergies:    allergies,
		Illnesses:    illnesses,
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
	member.Allergies = allergies
	member.Illnesses = illnesses
	member.Organization = organization
	member.BMI = utmath.BMICalculation(req.Weight, req.Height)

	return &member
}
