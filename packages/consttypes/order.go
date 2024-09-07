package consttypes

import (
	"fmt"
	"strings"
)

type (
	OrderStatus string
)

const (
	// * by member
	// * when member just checkout their cart
	OS_PLACED OrderStatus = "Placed"

	// * by partner
	// * when partner accept the order and update it
	OS_CONFIRMED      OrderStatus = "Confirmed"
	OS_BEING_PREPARED OrderStatus = "Being Prepared"
	OS_PREPARED       OrderStatus = "Prepared"
	OS_PICKED_UP      OrderStatus = "Picked Up"

	// * by member
	OS_COMPLETED OrderStatus = "Completed"

	// * automatically by system
	OS_CANCELLED OrderStatus = "Cancelled"
)

func (enum OrderStatus) String() string {
	return string(enum)
}

func NewOrderHistoryDescription(status OrderStatus, by string) string {
	status = OrderStatus(strings.ToLower(string(status)))

	switch status {
	case OS_PLACED:
		return fmt.Sprintf("Order was %s by %s .", status, by)
	case OS_CONFIRMED:
		return fmt.Sprintf("Order was %s by %s and is being processed.", status, by)
	case OS_BEING_PREPARED:
		return fmt.Sprintf("Order is currently %s by %s.", status, by)
	case OS_PREPARED:
		return fmt.Sprintf("Order has been %s by %s and is ready for pickup.", status, by)
	case OS_PICKED_UP:
		return fmt.Sprintf("Order was %s and is on its way.", status)
	case OS_COMPLETED:
		return fmt.Sprintf("Order has been marked as %s by %s. The customer has received the order.", status, by)
	case OS_CANCELLED:
		return fmt.Sprintf("Order was %s by %s. No further action will be taken.", status, by)
	default:
		return fmt.Sprintf("Order is %s by %s.", status, by)
	}
}
