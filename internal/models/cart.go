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
	Cart struct {
		base.Model

		MealID uuid.UUID `json:"meal_id" gorm:"required"`
		Meal   Meal      `json:"meal"`

		ReferenceID   uuid.UUID           `json:"reference_id" gorm:"required;<-:create"`
		ReferenceType consttypes.UserRole `json:"referenceType" gorm:"required; oneof='Member' 'Caregiver';type:user_role_enum;<-:create"`
		Quantity      uint                `json:"quantity" gorm:"required"`
	}
)

func (c *Cart) ToResponse(
	member *responses.Member,
	caregiver *responses.Caregiver,
) (*responses.Cart, error) {
	var cres responses.Cart

	if err := copier.CopyWithOption(&cres, &c, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	cres.Member = member
	cres.Caregiver = caregiver

	return &cres, nil
}
