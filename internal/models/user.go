package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/uttoken"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	User struct {
		helper.Model
		Address                []*Address          `json:"address,omitempty"`
		UserImage              *UserImage          `json:"user_image,omitempty"`
		Email                  string              `json:"email" gorm:"not null;unique" example:"email@email.com"`
		Password               string              `json:"-" gorm:"size:255;not null;" binding:"required" example:"password"`
		Role                   consttypes.UserRole `json:"role" gorm:"not null" example:"0" default:"0"`
		ResetPasswordToken     string              `json:"-"`
		ResetPasswordSentAt    time.Time           `json:"-"`
		ConfirmationToken      string              `json:"-"`
		ConfirmedAt            time.Time           `json:"confirmed_at"`
		ConfirmationSentAt     time.Time           `json:"-"`
		RefreshToken           string              `json:"-"`
		RefreshTokenExpiration string              `json:"-"`
	}

	UserImage struct {
		helper.Model
		UserID  uuid.UUID `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		ImageID uuid.UUID `json:"image_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
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

func (u *User) ToResponse() *responses.UserResponse {
	var ures responses.UserResponse
	err := copier.Copy(&ures, &u)
	if err != nil {
		utlogger.LogError(err)
	}

	return &ures
}

func (u *User) ToAuthResponse(token *uttoken.TokenHeader) *responses.AuthResponse {
	aures := responses.AuthResponse{
		Token:   token.AccessToken,
		Expires: token.AccessTokenExpires,
	}

	err := copier.Copy(&aures, &u)
	if err != nil {
		utlogger.LogError(err)
	}

	return &aures
}
