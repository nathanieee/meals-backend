package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateCart struct {
		MealID        uuid.UUID           `json:"meal_id"`
		ReferenceID   uuid.UUID           `json:"reference_id"`
		ReferenceType consttypes.UserRole `json:"reference_type"`
		Quantity      uint                `json:"quantity"`
	}

	UpdateCart struct {
		MealID        uuid.UUID           `json:"meal_id"`
		ReferenceID   uuid.UUID           `json:"reference_id"`
		ReferenceType consttypes.UserRole `json:"reference_type"`
		Quantity      uint                `json:"quantity"`
	}
)

func (req *CreateCart) ToModel() (*models.Cart, error) {
	var cart models.Cart

	if err := copier.CopyWithOption(&cart, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &cart, nil
}

func (req *UpdateCart) ToModel(
	cart models.Cart,
) (*models.Cart, error) {
	if err := copier.CopyWithOption(&cart, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &cart, nil
}
