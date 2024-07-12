package requests

import "project-skbackend/external/controllers/requests"

type (
	Geolocation struct {
		Longitude string `json:"longitude" form:"longitude" binding:"required"`
		Latitude  string `json:"latitude" form:"latitude" binding:"required"`
	}
)

func (req *Geolocation) ToDistanceMatrixGeolocation() *requests.Geolocation {
	return &requests.Geolocation{
		Lng: req.Longitude,
		Lat: req.Latitude,
	}
}
