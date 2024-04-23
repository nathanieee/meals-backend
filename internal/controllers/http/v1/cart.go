package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/models"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	cartroutes struct {
		cfg   *configs.Config
		scart cartservice.ICartService
		suser userservice.IUserService
	}
)

func newCartRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	scart cartservice.ICartService,
	suser userservice.IUserService,
) {
	r := &cartroutes{
		cfg:   cfg,
		scart: scart,
		suser: suser,
	}

	gcart := rg.Group("carts")
	{

		gmembcare := gcart.Group("").Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_MEMBER, consttypes.UR_CAREGIVER))
		{
			gmembcare.POST("", r.createCart)
		}

		gcart.GET("raw", r.getCartsRaw)
	}
}

// TODO - continue working on the create cart function
func (r *cartroutes) createCart(ctx *gin.Context) {
	var (
		function = "create cart"
		entity   = "cart"
		req      *requests.CreateCart
		err      error
	)

	userres, err := uttoken.GetUser(ctx)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	err = ctx.ShouldBindJSON(&req)
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

	roleres, err := r.suser.GetRoleDataByUserID(userres.ID)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	if roleres != nil {
		switch roleres.Role {
		case consttypes.UR_CAREGIVER:
			res, ok := roleres.Data.(*models.Caregiver)
			if !ok {
				utresponse.GeneralUnauthorized(
					ctx,
					consttypes.ErrUserInvalidRole,
				)
				return
			}

			req, err = req.New(res.ID, res.User.Role)
			if err != nil {
				utresponse.GeneralUnauthorized(
					ctx,
					err,
				)
				return
			}
		case consttypes.UR_MEMBER:
			res, ok := roleres.Data.(*models.Member)
			if !ok {
				utresponse.GeneralUnauthorized(
					ctx,
					consttypes.ErrUserInvalidRole,
				)
				return
			}

			req, err = req.New(res.ID, res.User.Role)
			if err != nil {
				utresponse.GeneralUnauthorized(
					ctx,
					err,
				)
				return
			}

		default:
			utresponse.GeneralUnauthorized(
				ctx,
				consttypes.ErrUserInvalidRole,
			)
			return
		}

	}

	rescart, err := r.scart.Create(*req)
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
		rescart,
	)
}

func (r *cartroutes) getCartsRaw(ctx *gin.Context) {
	var (
		entity = "carts"
	)

	carts, err := r.scart.Read()
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
		carts,
	)
}
