package responses

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type (
	UserResponse struct {
		ID                     uuid.UUID `json:"id"`
		CreatedAt              time.Time
		UpdatedAt              time.Time
		Address                []*models.Address   `json:"address,omitempty" gorm:"foreignKey:UserID;references:ID"`
		UserImage              *models.UserImage   `json:"userImage,omitempty" gorm:"foreignKey:UserID;references:ID"`
		Email                  string              `json:"email" example:"email@email.com"`
		Password               string              `json:"-" example:"password"`
		Role                   consttypes.UserRole `json:"role" example:"0"`
		ResetPasswordToken     int                 `json:"-"`
		ResetPasswordSentAt    time.Time           `json:"-"`
		ConfirmationToken      int                 `json:"-"`
		ConfirmedAt            time.Time           `json:"confirmedAt"`
		ConfirmationSentAt     time.Time           `json:"-"`
		RefreshToken           string              `json:"-"`
		RefreshTokenExpiration string              `json:"-"`
	}
)

func (ures *UserResponse) IsEmpty() bool {
	return cmp.Equal(ures, UserResponse{})
}
