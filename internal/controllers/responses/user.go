package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		base.Model

		Address []*Address `json:"address,omitempty"`

		Image *UserImage `json:"image,omitempty"`

		Email    string              `json:"email"`
		Password string              `json:"-"`
		Role     consttypes.UserRole `json:"role"`

		ResetPasswordToken  int       `json:"-"`
		ResetPasswordSentAt time.Time `json:"-"`

		ConfirmationToken  int       `json:"-"`
		ConfirmedAt        time.Time `json:"confirmed_at"`
		ConfirmationSentAt time.Time `json:"-"`
	}

	UserImage struct {
		base.Model

		UserID  uuid.UUID `json:"-"`
		ImageID uuid.UUID `json:"-"`

		Image Image `json:"image"`
	}
)
