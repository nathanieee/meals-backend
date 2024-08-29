package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/mealcategoryservice"
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
		smcat mealcategoryservice.IMealCategoryService
	}
)

func newMealRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	smeal mealservice.IMealService,
	smcat mealcategoryservice.IMealCategoryService,
) {
	r := &mealroutes{
		cfg:   cfg,
		smeal: smeal,
		smcat: smcat,
	}

	gmealpub := rg.Group("meals")
	{
		gmealpub.GET("", r.findMeals)
		gmealpub.GET("raw", r.findMealsRaw)
		gmealpub.GET(":mid", r.getMeal)

		gmcatpub := gmealpub.Group("categories")
		{
			gmcatpub.GET("", r.findMealCategories)
			gmcatpub.GET("raw", r.findMealCategoriesRaw)
			gmcatpub.GET(":mcid", r.getMealCategories)
		}
	}
}

func (r *mealroutes) findMeals(ctx *gin.Context) {
	var (
		entity  = "meals"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	meals, err := r.smeal.FindAll(reqpage)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
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

func (r *mealroutes) findMealsRaw(ctx *gin.Context) {
	var (
		entity = "meals"
	)

	meals, err := r.smeal.Read()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
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

func (r *mealroutes) getMeal(ctx *gin.Context) {
	var (
		function = "get meal"
		entity   = "meal"
	)

	uuid, err := uuid.Parse(ctx.Param("mid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	resmeal, err := r.smeal.GetByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		resmeal,
	)
}

func (r *mealroutes) findMealCategories(ctx *gin.Context) {
	var (
		entity  = "meal categories"
		reqpage = utrequest.GeneratePaginationFromRequest(ctx)
	)

	mcats, err := r.smcat.FindAll(reqpage)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		mcats,
	)
}

func (r *mealroutes) findMealCategoriesRaw(ctx *gin.Context) {
	var (
		entity = "meal categories"
	)

	mcats, err := r.smcat.Read()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		mcats,
	)
}

func (r *mealroutes) getMealCategories(ctx *gin.Context) {
	var (
		function = "get meal category"
		entity   = "meal category"
	)

	uuid, err := uuid.Parse(ctx.Param("mcid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	resmc, err := r.smcat.GetByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		resmc,
	)
}
