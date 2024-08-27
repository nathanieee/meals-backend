package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
)

type (
	profileroutes struct {
		cfg   *configs.Config
		suser userservice.IUserService
	}
)

func newProfileRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	suser userservice.IUserService,
) {
	r := &profileroutes{
		cfg:   cfg,
		suser: suser,
	}

	gprofilepvt := rg.Group("profiles")
	gprofilepvt.Use(middlewares.JWTAuthMiddleware(
		cfg,
		consttypes.UR_MEMBER,
		consttypes.UR_CAREGIVER,
		consttypes.UR_ADMIN,
		consttypes.UR_ORGANIZATION,
		consttypes.UR_PARTNER,
		consttypes.UR_PATRON,
	))
	{
		gprofilepvt.GET("me", r.getOwnProfile)
	}
}

func (r *profileroutes) getOwnProfile(ctx *gin.Context) {
	var (
		function = "get own profile"
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

	utresponse.GeneralSuccessFetch(
		function,
		ctx,
		roleres.Data,
	)
}
