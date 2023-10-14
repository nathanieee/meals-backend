package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/di"
	"project-skbackend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(ge *gin.Engine, db *gorm.DB, cfg *configs.Config, di *di.DependencyInjection) {
	ge.Use(gin.Logger())
	ge.Use(gin.Recovery())
	ge.Use(middlewares.CORSMiddleware())

	ge.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "oks"})
	})

	h := ge.Group("api/v1")
	{
		newUserRoutes(h, db, cfg, di.UserService, di.MailService)
		newAuthRoutes(h, cfg, di.AuthService)
		newMemberRoutes(h, db, cfg, di.MemberService)
	}
}
