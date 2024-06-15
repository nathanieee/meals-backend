package models

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Donation struct {
		base.Model

		PatronID uuid.UUID `json:"-" gorm:"required"`

		Value  float64                   `json:"value" gorm:"required"`
		Status consttypes.DonationStatus `json:"status" gorm:"required; type:donation_status_enum"`
	}
)
