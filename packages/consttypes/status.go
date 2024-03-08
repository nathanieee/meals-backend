package consttypes

import "encoding/json"

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
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}

func (enum DonationStatus) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}
