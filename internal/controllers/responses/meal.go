package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	MealResponse struct {
		helper.Model
		MealImage     *MealImageResponse     `json:"meal_image,omitempty"`
		MealIllnesses *[]MealIllnessResponse `json:"meal_illness,omitempty"`
		MealAllergies *[]MealAllergyResponse `json:"meal_allergy,omitempty"`
		PartnerID     uuid.UUID              `json:"partner_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Partner       PartnerResponse        `json:"partner"`
		Name          string                 `json:"name" gorm:"not null" binding:"required" example:"Nasi Goyeng"`
		Status        consttypes.MealStatus  `json:"status" gorm:"not null; type:meal_status_enum" binding:"required" example:"Active"`
		Description   string                 `json:"description" gorm:"size:255" example:"This meal is made using chicken and egg."`
	}

	MealImageResponse struct {
		helper.Model
		MealID  uuid.UUID     `json:"meal_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		ImageID uuid.UUID     `json:"image_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Image   ImageResponse `json:"image"`
	}

	MealIllnessResponse struct {
		helper.Model
		MealID    uuid.UUID       `json:"meal_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		IllnessID uuid.UUID       `json:"illness_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Illness   IllnessResponse `json:"illness"`
	}

	MealAllergyResponse struct {
		helper.Model
		MealID    uuid.UUID       `json:"meal_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID uuid.UUID       `json:"allergy_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Allergy   AllergyResponse `json:"allergy"`
	}
)
