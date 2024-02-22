package models

import (
	"project-skbackend/internal/models/helper"
)

type (
	Illness struct {
		helper.Model
		Name        string `json:"name" gorm:"required" example:"Cold Sore"`
		Description string `json:"description" gorm:"required" example:"Infection with the herpes simplex virus around the border of the lips."`
	}
)

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
