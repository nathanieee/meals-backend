package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Partner struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User      `json:"user"`

		Name string `json:"name" gorm:"required" example:"McDonald's"`

		MealCategories []*MealCategory `json:"meal_categories,omitempty" gorm:"many2many:partner_meal_category_composites;"`
	}
)

func (p *Partner) ToResponse() (*responses.Partner, error) {
	pres := responses.Partner{}

	if err := copier.CopyWithOption(&pres, &p, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &pres, nil
}
