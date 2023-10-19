package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Meal struct {
		helper.Model
		MealImage   *MealImage            `json:"meal_image,omitempty"`
		MealIllness *MealIllness          `json:"meal_illness,omitempty"`
		MealAllergy *MealAllergy          `json:"meal_allergy,omitempty"`
		PartnerID   uuid.UUID             `json:"partner_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Partner     Partner               `json:"partner"`
		Name        string                `json:"name" gorm:"not null" binding:"required" example:"Nasi Goyeng"`
		Status      consttypes.MealStatus `json:"meal_status" gorm:"not null" binding:"required" example:"Active"`
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
