package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Partner struct {
		helper.Model
		UserID uuid.UUID `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User      `json:"user"`
		Name   string    `json:"name" gorm:"not null" binding:"required" example:"McDonald's"`
	}
)

func (p *Partner) ToResponse() *responses.Partner {
	pres := responses.Partner{}

	if err := copier.CopyWithOption(&pres, &p, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &pres
}
