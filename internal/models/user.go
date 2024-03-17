package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/uttoken"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	User struct {
		base.Model

		Address []*Address `json:"address,omitempty"`

		Image *UserImage `json:"image,omitempty"`

		Email    string              `json:"email" gorm:"required;unique" example:"email@email.com"`
		Password string              `json:"password" gorm:"size:255;required;" example:"password"`
		Role     consttypes.UserRole `json:"role" gorm:"required;type:user_role_enum" example:"0" default:"0"`

		ConfirmationToken  int       `json:"-"`
		ConfirmedAt        time.Time `json:"confirmed_at"`
		ConfirmationSentAt time.Time `json:"-"`

		ResetPasswordToken  string    `json:"-"`
		ResetPasswordSentAt time.Time `json:"-"`
	}

	UserImage struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`

		ImageID uuid.UUID `json:"image_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Image   Image     `json:"image"`
	}
)

func (u *User) ToResponse() (*responses.User, error) {
	var ures responses.User
	err := copier.CopyWithOption(&ures, &u, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &ures, nil
}

func (u *User) ToAuth(token *uttoken.TokenHeader) *responses.Auth {
	aures := responses.Auth{
		Token:   token.AccessToken,
		Expires: token.AccessTokenExpires,
	}

	err := copier.CopyWithOption(&aures, &u, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		utlogger.Error(err)
	}

	return &aures
}
