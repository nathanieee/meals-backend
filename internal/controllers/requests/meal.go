package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateMeal struct {
		*CreateImage
		IllnessID   []*uuid.UUID          `json:"illness_id" form:"illness_id" binding:"-"`
		AllergyID   []*uuid.UUID          `json:"allergy_id" form:"allergy_id" binding:"-"`
		PartnerID   uuid.UUID             `json:"partner_id" form:"partner_id" binding:"required"`
		Name        string                `json:"name" form:"name" binding:"required"`
		Status      consttypes.MealStatus `json:"status" form:"status" binding:"required"`
		Description string                `json:"description" form:"description" binding:"required"`
	}

	UpdateMeal struct {
		*UpdateImage
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
) (*models.Meal, error) {
	var meal models.Meal

	if err := copier.CopyWithOption(&meal, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	meal.Name = strings.Title(req.Name)
	meal.Images = images
	meal.Illnesses = illnesses
	meal.Allergies = allergies
	meal.Partner = partner

	return &meal, nil
}

func (req *UpdateMeal) ToModel(
	meal models.Meal,
	images []*models.MealImage,
	illnesses []*models.MealIllness,
	allergies []*models.MealAllergy,
	partner models.Partner,
) (*models.Meal, error) {
	if err := copier.CopyWithOption(&meal, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	meal.Name = strings.Title(req.Name)
	meal.Images = images
	meal.Illnesses = illnesses
	meal.Allergies = allergies
	meal.Partner = partner

	return &meal, nil
}
