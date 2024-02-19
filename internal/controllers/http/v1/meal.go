package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/mealservice"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	r := &mealroutes{
		cfg:   cfg,
		smeal: smeal,
	}

	admgrp := rg.Group("meals")
	{
		admgrp.POST("", r.createMeal)
		admgrp.GET("", r.getMeals)
		admgrp.GET("raw", r.getMealsRaw)
		admgrp.PUT("/:uuid", r.updateMeal)
		admgrp.DELETE("/:uuid", r.deleteMeal)
	}
}

func (r *mealroutes) createMeal(ctx *gin.Context) {
	var function = "create meal"
	var entity = "meal"
	var req requests.CreateMeal

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	meres, err := r.smeal.Create(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		meres,
	)
}

func (r *mealroutes) getMeals(ctx *gin.Context) {
	var entity = "meals"
	var paginationReq = utrequest.GeneratePaginationFromRequest(ctx)

	meals, err := r.smeal.FindAll(paginationReq)
	if err != nil {
		utresponse.GeneralNotFound(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		meals,
	)
}

func (r *mealroutes) getMealsRaw(ctx *gin.Context) {
	var entity = "meals"

	meals, err := r.smeal.Read()
	if err != nil {
		utresponse.GeneralNotFound(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		meals,
	)
}

func (r *mealroutes) updateMeal(ctx *gin.Context) {
	var function = "update meal"
	var entity = "meal"
	var req requests.UpdateMeal

	err := ctx.ShouldBind(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smeal.FindByID(uuid)
	if err != nil {
		utresponse.GeneralNotFound(
			entity,
			ctx,
			err,
		)
		return
	}

	meres, err := r.smeal.Update(uuid, req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		meres,
	)
}

func (r *mealroutes) deleteMeal(ctx *gin.Context) {
	var function = "delete meal"
	var entity = "meal"

	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smeal.FindByID(uuid)
	if err != nil {
		utresponse.GeneralNotFound(
			entity,
			ctx,
			err,
		)
		return
	}

	err = r.smeal.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}
