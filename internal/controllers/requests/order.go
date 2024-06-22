package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateOrder struct {
		Meals []CreateOrderMeal `json:"meals" form:"meals" binding:"required,dive"`
	}

	CreateOrderMeal struct {
		MealID uuid.UUID `json:"meal_id" form:"meal_id" binding:"required"`

		Quantity uint `json:"quantity" form:"quantity" binding:"required"`
	}
)

func (req *CreateOrder) ToModel(
	member models.Member,
	userorder models.User,
	meals []models.OrderMeal,
) (*models.Order, error) {
	var (
		order models.Order
	)

	if err := copier.CopyWithOption(&order, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	order.MemberID = member.ID
	order.Member = member
	order.Meals = meals
	order.Status = consttypes.OS_PROCESSED
	order.History = append(order.History, models.OrderHistory{
		UserID:      userorder.ID,
		User:        userorder,
		Status:      consttypes.OS_PROCESSED,
		Description: consttypes.NewOrderHistoryDescription(consttypes.OS_PROCESSED, userorder.Email),
	})

	return &order, nil
}

func (req *CreateOrderMeal) ToModel(meal models.Meal) (*models.OrderMeal, error) {
	var (
		omeal models.OrderMeal
	)

	if err := copier.CopyWithOption(&omeal, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	omeal.MealID = meal.ID
	omeal.Meal = meal

	return &omeal, nil
}
