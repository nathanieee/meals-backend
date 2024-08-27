package controllers

import (
	"project-skbackend/configs"
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
	cfg *configs.Config,
	smeal mealservice.IMealService,
) {
	r := &mealroutes{
		cfg:   cfg,
		smeal: smeal,
	}

	gmealpub := rg.Group("meals")
	{
		gmealpub.GET("", r.findMeals)
		gmealpub.GET("raw", r.findMealsRaw)
		gmealpub.GET(":mid", r.getMeal)
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
