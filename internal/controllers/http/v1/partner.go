package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/partnerservice"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	partnerroutes struct {
		cfg      *configs.Config
		spartner partnerservice.IPartnerService
	}
)

func newPartnerRoutes(
	rg *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	spartner partnerservice.IPartnerService,
) {
	// r := &partnerroutes{
	// 	cfg:      cfg,
	// 	spartner: spartner,
	// }

	// admgrp := rg.Group("partners")
}
