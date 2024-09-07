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
		Quantity int       `json:"quantity" form:"quantity" binding:"required"`
	}

	UpdateCart struct {
		MealID   uuid.UUID `json:"meal_id" form:"meal_id"`
		Quantity int       `json:"quantity" form:"quantity"`
	}
)

func (req *CreateCart) ToModel(member models.Member, meal models.Meal) (*models.Cart, error) {
	var (
		cart models.Cart
	)

	if err := copier.CopyWithOption(&cart, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cart.MemberID = member.ID
	cart.Member = member

	cart.PartnerID = meal.PartnerID
	cart.Partner = meal.Partner

	return &cart, nil
}

func (req *UpdateCart) ToModel(
	cart models.Cart,
) (*models.Cart, error) {
	cart.MealID = req.MealID

	return &cart, nil
}
