package controllers

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userRoutes struct {
	s   services.IUserService
	cfg *configs.Config
}

func newUserRoutes(h *gin.RouterGroup, db *gorm.DB, s services.IUserService, cfg *configs.Config) {
	r := &userRoutes{s: s, cfg: cfg}

	userRole := h.Group("users")
	{
		userRole.GET("", r.getUsers)
	}
}

func (r *userRoutes) getUsers(ctx *gin.Context) {
	fmt.Println("test")
}
