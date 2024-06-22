package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Order struct {
		base.Model

		MemberID uuid.UUID `json:"member_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Member   Member    `json:"member"`

		Meals []OrderMeal `json:"meal" gorm:"foreignKey:OrderID"`

		Status consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`

		History []OrderHistory `json:"history" gorm:"foreignKey:OrderID"`
	}

	OrderMeal struct {
		base.Model

		OrderID uuid.UUID `json:"order_id" gorm:"required" example:"1"`

		MealID uuid.UUID `json:"meal_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Meal   Meal      `json:"meal"`

		Quantity uint `json:"quantity" gorm:"required" example:"2"`
	}

	OrderHistory struct {
		base.Model

		OrderID uuid.UUID `json:"order_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`

		UserID uuid.UUID `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User      `json:"user"`

		Status      consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`
		Description string                 `json:"description" gorm:"required" example:"This ores is made using chicken and egg."`
	}
)

func (o *Order) ToResponse() (*responses.Order, error) {
	var (
		ores responses.Order
	)

	if err := copier.CopyWithOption(&ores, &o, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &ores, nil
}
