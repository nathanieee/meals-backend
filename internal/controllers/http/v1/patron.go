package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/patronservice"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	patronroutes struct {
		cfg     *configs.Config
		spatron patronservice.IPatronService
		sauth   authservice.IAuthService
	}
)

func newPatronRoutes(

	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
	spatron patronservice.IPatronService,
) {
	r := &patronroutes{
		cfg:     cfg,
		sauth:   sauth,
		spatron: spatron,
	}

	gpatronspub := rg.Group("patrons")
	{
		gpatronspub.POST("register", r.patronRegister)
	}
}

func (r *patronroutes) patronRegister(ctx *gin.Context) {
	var (
		function = "patron register"
		entity   = "patron"
		req      requests.CreatePatron
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

	_, err = r.spatron.Create(req)
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
