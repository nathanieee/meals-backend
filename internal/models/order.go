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

		PartnerID uuid.UUID `json:"partner_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Partner   Partner   `json:"partner"`

		Meals []OrderMeal `json:"meals" gorm:"foreignKey:OrderID"`

		Status consttypes.OrderStatus `json:"status" gorm:"required; type:order_status_enum" example:"Pending"`

		History []OrderHistory `json:"histories" gorm:"foreignKey:OrderID"`
	}

	OrderMeal struct {
		base.Model

		OrderID uuid.UUID `json:"order_id" gorm:"required" example:"1"`

		MealID uuid.UUID `json:"meal_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Meal   Meal      `json:"meal"`

		PartnerID uuid.UUID `json:"partner_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Partner   Partner   `json:"partner"`

		Quantity int `json:"quantity" gorm:"required" example:"2"`
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

func NewCreateOrderMeals(
	meal Meal,
	quantity int,
) *OrderMeal {
	return &OrderMeal{
		MealID:    meal.ID,
		PartnerID: meal.PartnerID,
		Quantity:  quantity,
	}
}

func NewOrderHistory(
	user User,
	status consttypes.OrderStatus,
) *OrderHistory {
	return &OrderHistory{
		UserID:      user.ID,
		User:        user,
		Status:      status,
		Description: consttypes.NewOrderHistoryDescription(status, user.Email),
	}
}

func (o *Order) OrderConfirmed(user User) *Order {
	var (
		status = consttypes.OS_CONFIRMED
	)

	o.Status = status

	// * create new order history
	oh := NewOrderHistory(user, status)

	// * append history
	o.History = append(o.History, *oh)
	return o
}

func (o *Order) OrderBeingPrepared(user User) *Order {
	var (
		status = consttypes.OS_BEING_PREPARED
	)

	o.Status = status

	// * create new order history
	oh := NewOrderHistory(user, status)

	// * append history
	o.History = append(o.History, *oh)
	return o
}

func (o *Order) OrderPrepared(user User) *Order {
	var (
		status = consttypes.OS_PREPARED
	)

	o.Status = status

	// * create new order history
	oh := NewOrderHistory(user, status)

	// * append history
	o.History = append(o.History, *oh)
	return o
}

func (o *Order) OrderPickedUp(user User) *Order {
	var (
		status = consttypes.OS_PICKED_UP
	)

	o.Status = status

	// * create new order history
	oh := NewOrderHistory(user, status)

	// * append history
	o.History = append(o.History, *oh)
	return o
}
