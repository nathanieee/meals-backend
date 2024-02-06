package requests

type (
	CreateAddress struct {
		Name      string `json:"name" form:"name" binding:"required"`
		Address   string `json:"address" form:"address" binding:"required"`
		Note      string `json:"note" form:"note" binding:"required;max:32"`
		Longitude string `json:"longitude" form:"longitude" binding:"required,longitude"`
		Latitude  string `json:"latitude" form:"latitude" binding:"required,latitude"`
	}

	UpdateAddress struct {
		Name      string `json:"name" form:"name" binding:"required"`
		Address   string `json:"address" form:"address" binding:"required"`
		Note      string `json:"note" form:"note" binding:"required"`
		Longitude string `json:"longitude" form:"longitude" binding:"required,longitude"`
		Latitude  string `json:"latitude" form:"latitude" binding:"required,latitude"`
	}
)