package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Donation struct {
		base.Model

		Value  float64                   `json:"value"`
		Status consttypes.DonationStatus `json:"status"`
	}
)
