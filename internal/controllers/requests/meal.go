package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateMeal struct {
		MealImage   *CreateImage          `json:"meal_image" form:"meal_image" binding:"-"`
		IllnessID   []*uuid.UUID          `json:"illness_id" form:"illness_id" binding:"-"`
		AllergyID   []*uuid.UUID          `json:"allergy_id" form:"allergy_id" binding:"-"`
		PartnerID   uuid.UUID             `json:"partner_id" form:"partner_id" binding:"required"`
		Name        string                `json:"name" form:"name" binding:"required"`
		Status      consttypes.MealStatus `json:"status" form:"status" binding:"required"`
		Description string                `json:"description" form:"description" binding:"required"`
	}

	UpdateMeal struct {
		MealImage   *UpdateImage          `json:"meal_image" form:"meal_image" binding:"-"`
		IllnessID   []*uuid.UUID          `json:"illness_id" form:"illness_id" binding:"-"`
		AllergyID   []*uuid.UUID          `json:"allergy_id" form:"allergy_id" binding:"-"`
		PartnerID   uuid.UUID             `json:"partner_id" form:"partner_id" binding:"required"`
		Name        string                `json:"name" form:"name" binding:"required"`
		Status      consttypes.MealStatus `json:"status" form:"status" binding:"required"`
		Description string                `json:"description" form:"description" binding:"required"`
	}
)

func (req *CreateMeal) ToModel(
	images []*models.MealImage,
	illnesses []*models.MealIllness,
	allergies []*models.MealAllergy,
	partner models.Partner,
) *models.Meal {
	meal := models.Meal{
		Images:    images,
		Illnesses: illnesses,
		Allergies: allergies,
		Partner:   partner,
	}

	if err := copier.Copy(&meal, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &meal
}
