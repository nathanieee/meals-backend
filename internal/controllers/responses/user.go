package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	UserResponse struct {
		helper.Model
		Address                []*AddressResponse  `json:"address,omitempty"`
		UserImage              *UserImageResponse  `json:"user_image,omitempty"`
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

	UserImageResponse struct {
		helper.Model
		UserID  uuid.UUID     `json:"user_id"`
		ImageID uuid.UUID     `json:"image_id"`
		Image   ImageResponse `json:"image"`
	}
)
