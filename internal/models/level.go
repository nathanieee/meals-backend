package models

import "gorm.io/gorm"

type (
	Level struct {
		gorm.Model
		Name string `json:"name,omitempty"`
	}
)
