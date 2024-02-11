package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/mealservice"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	mealroutes struct {
		cfg   *configs.Config
		smeal mealservice.IMealService
	}
)

func newMealRoutes(
	rg *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	smeal mealservice.IMealService,
) {
	// r := &mealroutes{
	// 	cfg:   cfg,
	// 	smeal: smeal,
	// }

	// admgrp := rg.Group("meals")
	// {
	// 	admgrp.POST("", r.createMeal)
	// 	admgrp.GET("", r.getMeals)
	// 	admgrp.PUT("/:uuid", r.updateMeal)
	// 	admgrp.DELETE("/:uuid", r.deleteMeal)
	// }
}
