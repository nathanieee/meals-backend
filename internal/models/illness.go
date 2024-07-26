package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	Illness struct {
		base.Model

		Name        string `json:"name" gorm:"required" example:"Cold Sore"`
		Description string `json:"description" gorm:"required" example:"Infection with the herpes simplex virus around the border of the lips."`
	}
)

func (ill *Illness) ToResponse() (*responses.Illness, error) {
	var (
		illres = responses.Illness{}
	)

	if err := copier.CopyWithOption(&illres, &ill, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &illres, nil
}

func (ill *Illness) ToMemberIllness() *MemberIllness {
	return &MemberIllness{
		IllnessID: ill.ID,
		Illness:   *ill,
	}
}

func (ill *Illness) ToMealIllness() *MealIllness {
	return &MealIllness{
		IllnessID: ill.ID,
		Illness:   *ill,
	}
}
