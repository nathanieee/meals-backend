package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	User struct {
		helper.Model
		Address                []*Address          `json:"address,omitempty"`
		Password               string              `json:"-" gorm:"size:255;not null;" binding:"required" example:"password"`
		UserImage              *UserImage          `json:"userImage,omitempty"`
		Email                  string              `json:"email" gorm:"not null;unique" example:"email@email.com"`
		Role                   consttypes.UserRole `json:"role" gorm:"not null" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"0"`
		ResetPasswordToken     int                 `json:"-"`
		ResetPasswordSentAt    time.Time           `json:"-"`
		ConfirmationToken      int                 `json:"-"`
		ConfirmedAt            time.Time           `json:"confirmedAt"`
		ConfirmationSentAt     time.Time           `json:"-"`
		RefreshToken           string              `json:"-"`
		RefreshTokenExpiration string              `json:"-"`
	}

	UserImage struct {
		helper.Model
		UserID  uuid.UUID `json:"userID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		ImageID uuid.UUID `json:"imageID" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Image   Image     `json:"image"`
	}
)

func (u *User) hashPasswordIfNeeded() error {
	hash, err := helper.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.hashPasswordIfNeeded()
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.hashPasswordIfNeeded()
}

func (u *User) IsEmpty() bool {
	return cmp.Equal(u, User{})
}
