package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Meal struct {
		helper.Model
		Images      []*MealImage          `json:"images,omitempty"`
		Illnesses   []*MealIllness        `json:"illnesses,omitempty"`
		Allergies   []*MealAllergy        `json:"allergies,omitempty"`
		Partner     Partner               `json:"partner"`
		Name        string                `json:"name" gorm:"not null" binding:"required" example:"Nasi Goyeng"`
		Status      consttypes.MealStatus `json:"status" gorm:"not null; type:meal_status_enum" binding:"required" example:"Active"`
		Description string                `json:"description" gorm:"size:255" example:"This meal is made using chicken and egg."`
	}

	MealImage struct {
		helper.Model
		MealID  uuid.UUID `json:"meal_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		ImageID uuid.UUID `json:"image_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Image   Image     `json:"image"`
	}

	MealIllness struct {
		helper.Model
		MealID    uuid.UUID `json:"meal_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		IllnessID uuid.UUID `json:"illness_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Illness   Illness   `json:"illness"`
	}

	MealAllergy struct {
		helper.Model
		MealID    uuid.UUID `json:"meal_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID uuid.UUID `json:"allergy_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Allergy   Allergy   `json:"allergy"`
	}
)
