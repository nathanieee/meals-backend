package models

type (
	Geolocation struct {
		Longitude string `json:"longitude" gorm:"required;longitude"`
		Latitude  string `json:"latitude" gorm:"required;latitude"`
	}
)
