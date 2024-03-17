package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/partnerservice"

	"github.com/gin-gonic/gin"
)

type (
	partnerroutes struct {
		cfg      *configs.Config
		spartner partnerservice.IPartnerService
	}
)

func newPartnerRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	spartner partnerservice.IPartnerService,
) {
	// r := &partnerroutes{
	// 	cfg:      cfg,
	// 	spartner: spartner,
	// }

	// admgrp := rg.Group("partners")
}
