package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
)

type (
	Member struct {
		helper.Model
		User         User              `json:"user,omitempty"`
		Caregiver    *Caregiver        `json:"caregiver,omitempty"`
		Organization *Organization     `json:"organization,omitempty"`
		Illnesses    []*MemberIllness  `json:"illnesses,omitempty"`
		Allergies    []*MemberAllergy  `json:"allergies,omitempty"`
		Height       float64           `json:"height,omitempty"`
		Weight       float64           `json:"weight,omitempty"`
		BMI          float64           `json:"bmi,omitempty"`
		FirstName    string            `json:"first_name,omitempty"`
		LastName     string            `json:"last_name,omitempty"`
		Gender       consttypes.Gender `json:"gender,omitempty"`
		DateOfBirth  customs.CDT_DATE  `json:"date_of_birth,omitempty"`
	}

	MemberIllness struct {
		helper.Model `json:"-"`
		Illness      Illness `json:"illness,omitempty"`
	}

	MemberAllergy struct {
		helper.Model `json:"-"`
		Allergy      Allergy `json:"allergy,omitempty"`
	}
)
