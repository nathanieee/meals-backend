package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
)

type (
	Member struct {
		helper.Model
		User         User              `json:"user"`
		Caregiver    *Caregiver        `json:"caregiver,omitempty"`
		Organization *Organization     `json:"organization,omitempty"`
		Illnesses    []*MemberIllness  `json:"illnesses,omitempty"`
		Allergies    []*MemberAllergy  `json:"allergies,omitempty"`
		Height       float64           `json:"height"`
		Weight       float64           `json:"weight"`
		BMI          float64           `json:"bmi"`
		FirstName    string            `json:"first_name"`
		LastName     string            `json:"last_name"`
		Gender       consttypes.Gender `json:"gender"`
		DateOfBirth  customs.CDT_DATE  `json:"date_of_birth"`
	}

	MemberIllness struct {
		helper.Model `json:"-"`
		Illness      Illness `json:"illness"`
	}

	MemberAllergy struct {
		helper.Model `json:"-"`
		Allergy      Allergy `json:"allergy"`
	}
)
