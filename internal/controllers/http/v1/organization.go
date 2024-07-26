package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/organizationservice"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	organizationroutes struct {
		cfg   *configs.Config
		sorg  organizationservice.IOrganizationService
		sauth authservice.IAuthService
	}
)

func newOrganizationRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
	sorg organizationservice.IOrganizationService,
) {
	r := &organizationroutes{
		cfg:   cfg,
		sauth: sauth,
		sorg:  sorg,
	}

	gorganizationspub := rg.Group("organizations")
	{
		gorganizationspub.POST("register", r.organizationRegister)
	}
}

func (r *organizationroutes) organizationRegister(ctx *gin.Context) {
	var (
		function = "organization register"
		entity   = "organization"
		req      requests.CreateOrganization
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

	// * get the user image
	file := req.User.CreateImage
	if file != nil {
		// * check if the file is an image
		err := file.IsImage()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	_, err = r.sorg.Create(req)

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
	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resauth,
	)
}
