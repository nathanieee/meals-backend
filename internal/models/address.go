package models

import (
	"project-skbackend/internal/models/base"

	"github.com/google/uuid"
)

type (
	Address struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required"`
		User   User      `json:"user"`

		Name    string `json:"name" gorm:"required"`
		Address string `json:"address" gorm:"required"`
		Note    string `json:"note" gorm:"required;max:256"`

		AddressDetailID *uuid.UUID     `json:"address_detail_id"`
		AddressDetail   *AddressDetail `json:"address_detail"`
	}

	AddressDetail struct {
		base.Model

		Geolocation

		FormattedAddress string `json:"formatted_address" gorm:"default:null"`
		PostCode         string `json:"post_code" gorm:"default:null"`
		Country          string `json:"country" gorm:"default:null"`
	}
)
