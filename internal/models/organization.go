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
	Organization struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User      `json:"user"`

		Type consttypes.OrganizationType `json:"type" gorm:"required; type:organization_type_enum" example:"Nursing Home"`
		Name string                      `json:"name" gorm:"required" example:"Panti Jompo Syailendra"`
	}
)

func (o *Organization) ToResponse() (*responses.Organization, error) {
	var (
		orgres responses.Organization
	)

	if err := copier.CopyWithOption(&orgres, &o, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &orgres, nil
}
