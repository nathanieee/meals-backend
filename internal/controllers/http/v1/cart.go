package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
)

type (
	cartroutes struct {
		cfg   *configs.Config
		scart cartservice.ICartService
	}
)

func newCartRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	scart cartservice.ICartService,
) {
	r := &cartroutes{
		cfg:   cfg,
		scart: scart,
	}

	grp := rg.Group("carts")
	{
		grp.GET("raw", r.getCartsRaw)
	}
}

func (r *cartroutes) getCartsRaw(ctx *gin.Context) {
	var entity = "carts"

	carts, err := r.scart.Read()
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
		carts,
	)
}
