package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
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

	gcartspvt := rg.Group("carts")
	gcartspvt.Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_MEMBER, consttypes.UR_CAREGIVER))
	{
		gcartspvt.GET("own", r.getOwnCart)
	}
}

func (r *cartroutes) getOwnCart(ctx *gin.Context) {
	var (
		function = "get cart"
		entity   = "cart"
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

	roleres, err := r.suser.GetRoleDataByUserID(userres.ID)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	if roleres == nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			consttypes.ErrUserInvalidRole,
		)
		return
	}

	rescart, err := r.scart.FindByMemberID(*roleres)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		rescart,
	)
}
