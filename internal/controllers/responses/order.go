package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Order struct {
		base.Model

		Member Member `json:"member"`

		Partner Partner `json:"partner"`

		Meals []OrderMeal `json:"meals" gorm:"foreignKey:OrderID"`

		Status consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`

		History []OrderHistory `json:"histories" gorm:"foreignKey:OrderID"`
	}

	OrderMeal struct {
		base.Model

		Meal Meal `json:"meal"`

		Partner Partner `json:"partner"`

		Quantity uint `json:"quantity" gorm:"required" example:"2"`
	}

	OrderHistory struct {
		base.Model

		User User `json:"user"`

		Status      consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`
		Description string                 `json:"description" gorm:"required" example:"This ores is made using chicken and egg."`
	}

	OrderRemaining struct {
		Quantity int `json:"quantity" example:"2"`
	}
)

func NewOrderRemaining(quantity int) (*OrderRemaining, error) {
	return &OrderRemaining{
		Quantity: quantity,
	}, nil
}
