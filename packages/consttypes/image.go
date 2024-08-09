package consttypes

type (
	ImageType string
)

const (
	IT_PROFILE       ImageType = "Profile"
	IT_MEAL          ImageType = "Meal"
	IT_MEAL_CATEGORY ImageType = "Meal Category"
)

func (enum ImageType) String() string {
	return string(enum)
}
