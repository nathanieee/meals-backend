package responses

import (
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/go-cmp/cmp"
)

type (
	UserResponse struct {
		helper.Model
		Address                []*models.Address   `json:"address,omitempty" gorm:"foreignKey:UserID;references:ID"`
		UserImage              *models.UserImage   `json:"user_image,omitempty" gorm:"foreignKey:UserID;references:ID"`
		Email                  string              `json:"email" example:"email@email.com"`
		Password               string              `json:"-" example:"password"`
		Role                   consttypes.UserRole `json:"role" example:"0"`
		ResetPasswordToken     int                 `json:"-"`
		ResetPasswordSentAt    time.Time           `json:"-"`
		ConfirmationToken      int                 `json:"-"`
		ConfirmedAt            time.Time           `json:"confirmed_at"`
		ConfirmationSentAt     time.Time           `json:"-"`
		RefreshToken           string              `json:"-"`
		RefreshTokenExpiration string              `json:"-"`
	}
)

func (ures *UserResponse) IsEmpty() bool {
	return cmp.Equal(ures, UserResponse{})
}
