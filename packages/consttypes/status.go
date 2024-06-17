package consttypes

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

	OS_PENDING    OrderStatus = "Pending"
	OS_PREPARING  OrderStatus = "Preparing"
	OS_PREPARED   OrderStatus = "Prepared"
	OS_DELIVERING OrderStatus = "Delivering"
	OS_DELIVERED  OrderStatus = "Delivered"
	OS_CANCELED   OrderStatus = "Canceled"
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
