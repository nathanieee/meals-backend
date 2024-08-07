package responses

import "project-skbackend/internal/models/base"

type (
	Address struct {
		base.Model

		Name    string `json:"name"`
		Address string `json:"address"`
		Note    string `json:"note"`

		AddressDetail AddressDetail `json:"address_detail"`
	}

	AddressDetail struct {
		base.Model

		Geolocation

		FormattedAddress string `json:"formatted_address"`
		PostCode         string `json:"post_code"`
		Country          string `json:"country"`
	}
)
