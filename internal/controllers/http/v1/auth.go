package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
)

type (
	authroutes struct {
		cfg   *configs.Config
		sauth authservice.IAuthService
	}
)

func newAuthRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
) {
	r := &authroutes{
		cfg:   cfg,
		sauth: sauth,
	}

	guser := rg.Group("auth")
	{
		guser.POST("signin", r.signin)

		// TODO - add a verify route group to verify email

		guser.POST("forgot-password", r.forgotPassword)
		guser.POST("reset-password", r.resetPassword)
		guser.GET("refresh-token", r.refreshAuthToken)
	}
}

func (r *authroutes) signin(
	ctx *gin.Context,
) {
	var function = "signin"
	var req requests.Signin

	err := ctx.ShouldBindJSON(&req)
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

	resuser, thead, err := r.sauth.Signin(req, ctx)
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

func (r *authroutes) forgotPassword(
	ctx *gin.Context,
) {
	var function = "forgot password"
	var req requests.ForgotPassword

	err := ctx.ShouldBindJSON(&req)
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

	err = r.sauth.ForgotPassword(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccess(
		function,
		ctx,
		nil,
	)
}

func (r *authroutes) resetPassword(
	ctx *gin.Context,
) {
	var function = "reset password"
	var req requests.ResetPassword

	err := ctx.ShouldBindJSON(&req)
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

	err = r.sauth.ResetPassword(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccess(
		function,
		ctx,
		nil,
	)
}

func (r *authroutes) refreshAuthToken(
	ctx *gin.Context,
) {
	var function = "refresh token"

	refreshToken, err := ctx.Cookie("refresh-token")
	if refreshToken == "" || err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			utresponse.ErrTokenNotFound,
		)
		return
	}

	resuser, thead, err := r.sauth.RefreshAuthToken(refreshToken, ctx)
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
