package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateCart struct {
		MealID   uuid.UUID `json:"meal_id" form:"meal_id" binding:"required"`
		Quantity uint      `json:"quantity" form:"quantity" binding:"required"`
	}

	UpdateCart struct {
		MealID   uuid.UUID `json:"meal_id" form:"meal_id" binding:"required"`
		Quantity uint      `json:"quantity" form:"quantity" binding:"required"`
	}
)

func (req *CreateCart) ToModel(member models.Member) (*models.Cart, error) {
	var (
		cart models.Cart
	)

	if err := copier.CopyWithOption(&cart, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cart.MemberID = member.ID
	cart.Member = member

	return &cart, nil
}

func (req *UpdateCart) ToModel(
	cart models.Cart,
) (*models.Cart, error) {
	if err := copier.CopyWithOption(&cart, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &cart, nil
}
