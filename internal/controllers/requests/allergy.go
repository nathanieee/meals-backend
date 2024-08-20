package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreateAllergy struct {
		Name        string               `json:"name" form:"name" binding:"required"`
		Description string               `json:"description" form:"description" binding:"required"`
		Allergens   consttypes.Allergens `json:"allergens" form:"allergens" binding:"required"`
	}

	UpdateAllergy struct {
		Name        string               `json:"name" form:"name"`
		Description string               `json:"description" form:"description"`
		Allergens   consttypes.Allergens `json:"allergens" form:"allergens"`
	}
)

func (req *CreateAllergy) ToModel() *models.Allergy {
	return &models.Allergy{
		Name:        req.Name,
		Description: req.Description,
		Allergens:   req.Allergens,
	}
}

func (req *UpdateAllergy) ToModel(all models.Allergy) (*models.Allergy, error) {
	if req == nil {
		return &all, nil
	}

	if err := copier.CopyWithOption(&all, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &all, nil
}
