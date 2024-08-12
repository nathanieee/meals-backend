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
		CartIDs []uuid.UUID `json:"cart_ids" form:"cart_ids" binding:"required"`
	}
)

func (req *CreateOrder) ToModel(
	member models.Member,
	userorder models.User,
	meals []models.OrderMeal,
) (*models.Order, error) {
	var (
		order  models.Order
		status = consttypes.OS_PLACED
	)

	if err := copier.CopyWithOption(&order, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	order.MemberID = member.ID
	order.Member = member
	order.Meals = meals
	order.Status = status

	// * create new order history
	oh := models.NewOrderHistory(userorder, status)

	// * append history
	order.History = append(order.History, *oh)

	return &order, nil
}
