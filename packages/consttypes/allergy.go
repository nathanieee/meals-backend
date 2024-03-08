package consttypes

import "encoding/json"

type (
	Allergens string
)

const (
	A_FOOD          Allergens = "Food"
	A_MEDICAL       Allergens = "Medical"
	A_ENVIRONMENTAL Allergens = "Environmental"
	A_CONTACT       Allergens = "Contact"
)

func (enum Allergens) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}
