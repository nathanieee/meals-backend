package consttypes

import (
	"fmt"
)

type (
	MealStatus     string
	DonationStatus string
	OrderStatus    string
)

const (
	MS_ACTIVE     MealStatus = "Active"
	MS_INACTIVE   MealStatus = "Inactive"
	MS_OUTOFSTOCK MealStatus = "Out of Stock"

	DS_ACCEPTED DonationStatus = "Accepted"
	DS_REJECTED DonationStatus = "Rejected"

	OS_PROCESSED OrderStatus = "processed"
	OS_PREPARED  OrderStatus = "prepared"
	OS_DELIVERED OrderStatus = "delivered"
	OS_CANCELED  OrderStatus = "canceled"
)

func (enum MealStatus) String() string {
	return string(enum)
}

func (enum DonationStatus) String() string {
	return string(enum)
}

func (enum OrderStatus) String() string {
	return string(enum)
}

func NewOrderHistoryDescription(status OrderStatus, by string) string {
	return fmt.Sprintf("Order is %s by %s.", status, by)
}
