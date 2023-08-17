package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Password               string    `json:"-" gorm:"size:255;not null;" binding:"required" example:"password"`
		FullName               string    `json:"fullName" gorm:"not null" example:"user name"`
		Email                  string    `json:"email" gorm:"not null;unique" example:"email@email.com"`
		RoleID                 uint      `json:"roleID" gorm:"not null" example:"1"`
		Role                   Role      `json:"-"`
		ResetPasswordToken     string    `json:"-"`
		ResetPasswordSentAt    time.Time `json:"-"`
		ConfirmationToken      int       `json:"-"`
		ConfirmedAt            time.Time `json:"confirmedAt"`
		ConfirmationSentAt     time.Time `json:"-"`
		RefreshToken           string    `json:"-"`
		RefreshTokenExpiration string    `json:"-"`
	}
)
