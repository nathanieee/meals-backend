package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Donation struct {
		base.Model

		PatronID uuid.UUID `json:"-" gorm:"not null"`

		// * One-to-One relation with DonationProof
		Proof *DonationProof `json:"proof,omitempty" gorm:"foreignKey:DonationID;constraint:OnDelete:CASCADE;"`

		Value  float64                   `json:"value" gorm:"not null"`
		Status consttypes.DonationStatus `json:"status" gorm:"not null;type:donation_status_enum"`
	}

	DonationProof struct {
		base.Model

		// * Foreign key to Donation
		DonationID uuid.UUID `json:"donation_id" gorm:"not null"`

		ImageID uuid.UUID `json:"image_id" gorm:"not null"`
		Image   Image     `json:"image"`
	}
)

func (d *Donation) ToResponse() (*responses.Donation, error) {
	var (
		dres responses.Donation
	)

	if err := copier.CopyWithOption(&dres, &d, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &dres, nil
}
