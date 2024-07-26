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
	Patron struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required"`
		User   User      `json:"user"`

		Type consttypes.PatronType `json:"type" gorm:"required; type:patron_type_enum"`
		Name string                `json:"name" gorm:"required" example:"Anonymus"`

		Donations []Donation `json:"donations"`
	}
)

func (p *Patron) ToResponse() (*responses.Patron, error) {
	pres := responses.Patron{}

	if err := copier.CopyWithOption(&pres, &p, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &pres, nil
}
