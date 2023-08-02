package models

import "gorm.io/gorm"

type (
	Level struct {
		gorm.Model
		ID   uint   `json:"id" gorm:"primary_key" example:"999"`
		Name string `json:"name,omitempty"`
	}

	Role struct {
		gorm.Model
		ID      uint   `json:"id" gorm:"primary_key" example:"999"`
		Name    string `json:"name,omitempty"`
		LevelID uint   `json:"levelID,omitempty"`
		Level   Level  `json:"-"`
	}
)
