package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		helper.Model
		Address                []*Address          `json:"address,omitempty"`
		UserImage              *UserImage          `json:"user_image,omitempty"`
		Email                  string              `json:"email"`
		Password               string              `json:"-"`
		Role                   consttypes.UserRole `json:"role"`
		ResetPasswordToken     int                 `json:"-"`
		ResetPasswordSentAt    time.Time           `json:"-"`
		ConfirmationToken      int                 `json:"-"`
		ConfirmedAt            time.Time           `json:"confirmed_at,omitempty"`
		ConfirmationSentAt     time.Time           `json:"-"`
		RefreshToken           string              `json:"-"`
		RefreshTokenExpiration string              `json:"-"`
	}

	UserImage struct {
		helper.Model
		UserID  uuid.UUID `json:"user_id"`
		ImageID uuid.UUID `json:"image_id"`
		Image   Image     `json:"image"`
	}
)
