package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Cart struct {
		base.Model

		MealID uuid.UUID `json:"meal_id" gorm:"required"`
		Meal   Meal      `json:"meal"`

		MemberID uuid.UUID `json:"member_id" gorm:"required"`
		Member   Member    `json:"member"`

		Quantity uint `json:"quantity" gorm:"required"`
	}
)

func (c *Cart) ToResponse(
	member *responses.Member,
) (*responses.Cart, error) {
	var (
		cres responses.Cart
	)

	if err := copier.CopyWithOption(&cres, &c, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cres.Member = member

	return &cres, nil
}
