package responses

import (
	"project-skbackend/packages/consttypes"
	"time"

	"github.com/google/uuid"
)

type (
	Auth struct {
		ID                 uuid.UUID           `json:"id"`
		Email              string              `json:"email"`
		Role               consttypes.UserRole `json:"role"`
		ConfirmationSentAt time.Time           `json:"confirmation_sent_at,omitempty"`
		ConfirmedAt        time.Time           `json:"confirmed_at,omitempty"`
		CreatedAt          time.Time           `json:"created_at,omitempty"`
		UpdatedAt          time.Time           `json:"updated_at,omitempty"`
		Token              string              `json:"token,omitempty"`
		Expires            time.Time           `json:"expires,omitempty"`
	}
)
