package consttypes

type (
	MealStatus string
)

const (
	MS_ACTIVE     MealStatus = "Active"
	MS_INACTIVE   MealStatus = "Inactive"
	MS_OUTOFSTOCK MealStatus = "Out of Stock"
)

func (enum MealStatus) String() string {
	return string(enum)
}
