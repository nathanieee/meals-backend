package consttypes

type (
	MealStatus     string
	DonationStatus string
)

const (
	MS_ACTIVE     MealStatus = "Active"
	MS_INACTIVE   MealStatus = "Inactive"
	MS_OUTOFSTOCK MealStatus = "Out of Stock"

	DS_ACCEPTED DonationStatus = "Accepted"
	DS_REJECTED DonationStatus = "Rejected"
)

func (enum MealStatus) String() string {
	return string(enum)
}

func (enum DonationStatus) String() string {
	return string(enum)
}
