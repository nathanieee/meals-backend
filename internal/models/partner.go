package models

import (
	"fmt"
	"project-skbackend/internal/models/helper"

	"gorm.io/gorm"
)

type (
	Partner struct {
		helper.Model
		UserID int    `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User   `json:"user"`
		Name   string `json:"name" gorm:"not null" binding:"required" example:"McDonald's"`
	}
)

func (p *Partner) BeforeCreate(tx *gorm.DB) error {
	var user *User

	if err := tx.Where("email = ?", p.User.Email).First(&user).Error; err != nil {
		return err
	}

	if user != nil {
		var err = fmt.Errorf("user with email %s already exists", p.User.Email)
		return err
	}

	return nil
}
