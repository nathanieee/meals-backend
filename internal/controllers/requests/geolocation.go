package requests

import "project-skbackend/external/controllers/requests"

type (
	Geolocation struct {
		Longitude string `json:"longitude" form:"longitude" binding:"required"`
		Latitude  string `json:"latitude" form:"latitude" binding:"required"`
	}

	DistanceMatrix struct {
		Origins      Geolocation `json:"origins"`
		Destinations Geolocation `json:"destinations"`
	}
)

func (req *Geolocation) ToDistanceMatrixGeolocation() *requests.Geolocation {
	return &requests.Geolocation{
		Lng: req.Longitude,
		Lat: req.Latitude,
	}
}

func (req *DistanceMatrix) ToDistanceMatrix() *requests.DistanceMatrix {
	return &requests.DistanceMatrix{
		Origins:      *req.Origins.ToDistanceMatrixGeolocation(),
		Destinations: *req.Destinations.ToDistanceMatrixGeolocation(),
	}
}
