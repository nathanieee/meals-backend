package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Donation struct {
		helper.Model
		Value  float64                   `json:"value" gorm:"required"`
		Status consttypes.DonationStatus `json:"status" gorm:"required; type:donation_status_enum"`
	}
)
