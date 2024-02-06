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
		Illness        []*MemberIllness  `json:"illnesses,omitempty"`
		Allergy        []*MemberAllergy  `json:"allergies,omitempty"`
		Height         float64           `json:"height"`
		Weight         float64           `json:"weight"`
		BMI            float64           `json:"bmi"`
		FirstName      string            `json:"first_name"`
		LastName       string            `json:"last_name"`
		Gender         consttypes.Gender `json:"gender"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth"`
	}

	MemberIllness struct {
		helper.Model `json:"-"`
		IllnessID    uuid.UUID `json:"-"`
		Illness      Illness   `json:"illness"`
	}

	MemberAllergy struct {
		helper.Model `json:"-"`
		AllergyID    uuid.UUID `json:"-"`
		Allergy      Allergy   `json:"allergy"`
	}
)
