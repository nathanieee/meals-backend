package requests

type (
	Geolocation struct {
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	}

	DistanceMatrix struct {
		Origins      Geolocation `json:"origins"`
		Destinations Geolocation `json:"destinations"`
	}
)
