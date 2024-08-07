package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/di"
	"project-skbackend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(ge *gin.Engine, db *gorm.DB, cfg *configs.Config, di *di.DependencyInjection, rdb *redis.Client) {
	ge.Use(gin.Logger())
	ge.Use(gin.Recovery())
	ge.Use(middlewares.CORSMiddleware())

	ge.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "oks"})
	})

	h := ge.Group("api/v1")
	{
		newAuthRoutes(h, cfg, rdb, di.AuthService, di.UserService)
		newMemberRoutes(h, cfg, di.MemberService, di.CartService, di.UserService, di.AuthService, di.OrderService)
		newPartnerRoutes(h, cfg, di.AuthService, di.PartnerService)
		newManageRoutes(h, cfg, di.MealService, di.MemberService, di.PartnerService, di.PatronService, di.IllnessService)
		newPatronRoutes(h, cfg, di.AuthService, di.PatronService)
		newOrganizationRoutes(h, cfg, di.AuthService, di.OrganizationService)
		newFileRoutes(h, cfg, di.FileService)
	}
}
