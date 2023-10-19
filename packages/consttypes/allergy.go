package consttypes

type (
	Allergens string
)

const (
	A_FOOD          Allergens = "Food"
	A_MEDICAL       Allergens = "Medical"
	A_ENVIRONMENTAL Allergens = "Environmental"
	A_CONTACT       Allergens = "Contact"
)
