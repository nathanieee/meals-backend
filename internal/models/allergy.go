package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	Allergy struct {
		base.Model

		Name        string               `json:"name" gorm:"required" example:"Milk"`
		Description string               `json:"description" gorm:"required" example:"A milk allergy, also known as a dairy allergy, is an adverse immune system response to one or more proteins found in cow's milk. It is different from lactose intolerance, which is a non-immune digestive disorder where the body has difficulty digesting lactose, a sugar found in milk. A milk allergy is an immune system disorder and can be more severe."`
		Allergens   consttypes.Allergens `json:"allergens" gorm:"required; type:allergens_enum" example:"Food"`
	}
)

func (all *Allergy) ToResponse() (*responses.Allergy, error) {
	var (
		allres = responses.Allergy{}
	)

	if err := copier.CopyWithOption(&allres, &all, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &allres, nil
}

func (ally *Allergy) ToMemberAllergy() *MemberAllergy {
	return &MemberAllergy{
		AllergyID: ally.ID,
		Allergy:   *ally,
	}
}

func (ally *Allergy) ToMealAllergy() *MealAllergy {
	return &MealAllergy{
		AllergyID: ally.ID,
		Allergy:   *ally,
	}
}
