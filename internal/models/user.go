package models

import (
	"fmt"
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
		Image                  *UserImage          `json:"image,omitempty"`
		Email                  string              `json:"email" gorm:"not null;unique" example:"email@email.com"`
		Password               string              `json:"-" gorm:"size:255;not null;" binding:"required" example:"password"`
		Role                   consttypes.UserRole `json:"role" gorm:"not null" example:"0" default:"0"`
		ResetPasswordToken     string              `json:"-"`
		ResetPasswordSentAt    time.Time           `json:"-"`
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

func (u *User) checkDuplicateEmail(tx *gorm.DB) error {
	var user *User

	result := tx.Where("email = ?", u.Email).First(&user)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if result.RowsAffected != 0 {
		var err = fmt.Errorf("user with email %s already exists", u.Email)
		return err
	}

	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.checkDuplicateEmail(tx); err != nil {
		return err
	}

	if err := u.hashPasswordIfNeeded(); err != nil {
		return err
	}

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.hashPasswordIfNeeded()
}

func (u *User) ToResponse() *responses.User {
	var ures responses.User
	err := copier.Copy(&ures, &u)
	if err != nil {
		utlogger.LogError(err)
	}

	return &ures
}

func (u *User) ToAuth(token *uttoken.TokenHeader) *responses.Auth {
	aures := responses.Auth{
		Token:   token.AccessToken,
		Expires: token.AccessTokenExpires,
	}

	err := copier.Copy(&aures, &u)
	if err != nil {
		utlogger.LogError(err)
	}

	return &aures
}
