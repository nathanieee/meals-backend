package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Order struct {
		base.Model

		Member Member `json:"member"`

		Meals []OrderMeal `json:"meal" gorm:"foreignKey:OrderID"`

		Status consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`

		History []OrderHistory `json:"history" gorm:"foreignKey:OrderID"`
	}

	OrderMeal struct {
		base.Model

		Meal Meal `json:"meal"`

		Quantity uint `json:"quantity" gorm:"required" example:"2"`
	}

	OrderHistory struct {
		base.Model

		User User `json:"user"`

		Status      consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`
		Description string                 `json:"description" gorm:"required" example:"This ores is made using chicken and egg."`
	}
)
