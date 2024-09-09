package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs/ctdatatype"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Caregiver struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User      `json:"user"`

		Gender      consttypes.Gender   `json:"gender" gorm:"required; type:gender_enum" example:"Male"`
		FirstName   string              `json:"first_name" gorm:"required" example:"Jonathan"`
		LastName    string              `json:"last_name" gorm:"required" example:"Vince"`
		DateOfBirth ctdatatype.CDT_DATE `json:"date_of_birth" gorm:"required" example:"2000-12-30"`
	}
)

func (c *Caregiver) ToResponse() (*responses.Caregiver, error) {
	var (
		cres responses.Caregiver
	)

	if err := copier.CopyWithOption(&cres, &c, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &cres, nil
}
