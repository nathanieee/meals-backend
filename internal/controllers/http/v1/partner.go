package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/partnerservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	partnerroutes struct {
		cfg      *configs.Config
		sauth    authservice.IAuthService
		spartner partnerservice.IPartnerService
	}
)

func newPartnerRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
	spartner partnerservice.IPartnerService,
) {
	r := &partnerroutes{
		cfg:      cfg,
		sauth:    sauth,
		spartner: spartner,
	}

	gpartnerspub := rg.Group("partners")
	{
		gpartnerspub.POST("register", r.partnerRegister)
	}

	gpartnerspvt := rg.Group("partners")
	gpartnerspvt.Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_PARTNER))
	{

	}
}

func (r *partnerroutes) partnerRegister(ctx *gin.Context) {
	var (
		function = "partner register"
		req      requests.CreatePartner
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	_, err = r.spartner.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			utresponse.GeneralDuplicate(
				"email",
				ctx,
				err,
			)
		} else {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
		}
		return
	}

	resuser, thead, err := r.sauth.Signin(*req.ToSignin(), ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	err = r.sauth.SendVerificationEmail(resuser.ID)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	resauth := thead.ToAuthResponse(*resuser)
	utresponse.GeneralSuccessAuth(
		function,
		ctx,
		resauth,
		thead,
	)
}
