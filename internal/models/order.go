package models

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Order struct {
		base.Model

		MemberID uuid.UUID `json:"member_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Member   Member    `json:"member"`

		MealID uuid.UUID `json:"meal_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Meal   Meal      `json:"meal"`

		Status consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`
	}

	OrderHistory struct {
		base.Model

		OrderID uuid.UUID `json:"order_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Order   Order     `json:"order"`

		Status      consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`
		Description string                 `json:"description" gorm:"required" example:"This order is made using chicken and egg."`
	}
)
