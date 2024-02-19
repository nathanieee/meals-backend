package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Meal struct {
		helper.Model
		Images      []*MealImage          `json:"images,omitempty"`
		Illnesses   []*MealIllness        `json:"illnesses,omitempty"`
		Allergies   []*MealAllergy        `json:"allergies,omitempty"`
		PartnerID   uuid.UUID             `json:"partner_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Partner     Partner               `json:"partner"`
		Name        string                `json:"name" gorm:"not null" binding:"required" example:"Nasi Goyeng"`
		Status      consttypes.MealStatus `json:"status" gorm:"not null; type:meal_status_enum" binding:"required" example:"Active"`
		Description string                `json:"description" example:"This meal is made using chicken and egg."`
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

func (m *Meal) ToResponse() *responses.Meal {
	mres := responses.Meal{}

	if err := copier.CopyWithOption(&mres, &m, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &mres
}
