package requests

type (
	CreateAddress struct {
		Name        string  `json:"name" binding:"required"`
		Address     string  `json:"address" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Note        string  `json:"note" binding:"required"`
		Landmark    string  `json:"landmark" binding:"required"`
		Longitude   float64 `json:"langitude" binding:"required"`
		Latitude    float64 `json:"latitude" binding:"required"`
	}
)
