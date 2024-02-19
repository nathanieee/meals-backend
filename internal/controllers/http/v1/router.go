package controllers

import (
	"mime/multipart"
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/di"
	"project-skbackend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		newMealRoutes(h, db, cfg, di.MealService)
		newPartnerRoutes(h, db, cfg, di.PartnerService)
	}
}

func ValidateImage(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}

	// Check if the uploaded file is an image (you may need to improve this check)
	return file.Header.Get("Content-Type") == "image/jpeg" || file.Header.Get("Content-Type") == "image/png"
}
