package requests

type (
	CreateAddress struct {
		Name    string `json:"name" form:"name" binding:"required"`
		Address string `json:"address" form:"address" binding:"required"`
		Note    string `json:"note" form:"note" binding:"required,max:256"`

		Geolocation
	}

	UpdateAddress struct {
		Name    string `json:"name" form:"name" binding:"-"`
		Address string `json:"address" form:"address" binding:"-"`
		Note    string `json:"note" form:"note" binding:"max:256"`

		Geolocation
	}
)
