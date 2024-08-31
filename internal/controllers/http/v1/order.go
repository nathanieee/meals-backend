package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/orderservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
)

type (
	orderroutes struct {
		cfg   *configs.Config
		sordr orderservice.IOrderService
		suser userservice.IUserService
	}
)

func newOrderRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sordr orderservice.IOrderService,
	suser userservice.IUserService,
) {
	r := &orderroutes{
		cfg:   cfg,
		sordr: sordr,
		suser: suser,
	}

	gordrpvt := rg.Group("orders")
	gordrpvt.Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_MEMBER, consttypes.UR_CAREGIVER))
	{
		gordrpvt.GET("own", r.getOwnOrder)
	}
}

func (r *orderroutes) getOwnOrder(ctx *gin.Context) {
	var (
		function = "get order"
		entity   = "order"
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

	resorder, err := r.sordr.FindByRoleRes(*roleres)
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
		resorder,
	)
}
