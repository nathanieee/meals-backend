package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"

	"github.com/google/uuid"
)

type (
	Member struct {
		helper.Model
		UserID         uuid.UUID         `json:"-"`
		User           User              `json:"user"`
		CaregiverID    *uuid.UUID        `json:"-"`
		Caregiver      *Caregiver        `json:"caregiver,omitempty"`
		OrganizationID *uuid.UUID        `json:"-"`
		Organization   *Organization     `json:"organization,omitempty"`
		Illness        []*MemberIllness  `json:"illness,omitempty"`
		Allergy        []*MemberAllergy  `json:"allergy,omitempty"`
		Height         float64           `json:"height"`
		Weight         float64           `json:"weight"`
		BMI            float64           `json:"bmi"`
		FirstName      string            `json:"first_name"`
		LastName       string            `json:"last_name"`
		Gender         consttypes.Gender `json:"gender"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth"`
	}

	MemberIllness struct {
		helper.Model
		IllnessID uuid.UUID `json:"illness_id"`
		Illness   Illness   `json:"illness"`
	}

	MemberAllergy struct {
		helper.Model
		AllergyID uuid.UUID `json:"allergy_id"`
		Allergy   Allergy   `json:"allergy"`
	}
)
