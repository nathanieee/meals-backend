package models

import (
	"fmt"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Patron struct {
		helper.Model
		UserID uuid.UUID             `json:"user_id" gorm:"not null" binding:"required"`
		User   User                  `json:"user"`
		Type   consttypes.PatronType `json:"type" gorm:"not null; type:patron_type_enum" binding:"required"`
		Name   string                `json:"name" gorm:"not null" binding:"required" example:"Anonymus"`
	}
)

func (p *Patron) BeforeCreate(tx *gorm.DB) error {
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
