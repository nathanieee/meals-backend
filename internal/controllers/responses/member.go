package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/custom"

	"github.com/google/uuid"
)

type (
	MemberResponse struct {
		helper.Model
		UserID         uuid.UUID                `json:"-"`
		User           UserResponse             `json:"user"`
		CaregiverID    *uuid.UUID               `json:"-"`
		Caregiver      *CaregiverResponse       `json:"caregiver,omitempty"`
		OrganizationID *uuid.UUID               `json:"-"`
		Organization   *OrganizationResponse    `json:"organization,omitempty"`
		Illness        []*MemberIllnessResponse `json:"illness,omitempty"`
		Allergy        []*MemberAllergyResponse `json:"allergy,omitempty"`
		Height         float64                  `json:"height"`
		Weight         float64                  `json:"weight"`
		BMI            float64                  `json:"bmi"`
		FirstName      string                   `json:"first_name"`
		LastName       string                   `json:"last_name"`
		Gender         consttypes.Gender        `json:"gender"`
		DateOfBirth    custom.CDT_DATE          `json:"date_of_birth"`
	}

	MemberIllnessResponse struct {
		helper.Model
		IllnessID uuid.UUID       `json:"illness_id"`
		Illness   IllnessResponse `json:"illness"`
	}

	MemberAllergyResponse struct {
		helper.Model
		AllergyID uuid.UUID       `json:"allergy_id"`
		Allergy   AllergyResponse `json:"allergy"`
	}
)
