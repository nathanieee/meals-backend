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

	usergrp := rg.Group("auth")
	{
		usergrp.POST("login", r.login)
		usergrp.POST("register", r.register)

		usergrp.POST("forgot-password", r.forgotPassword)
		usergrp.POST("reset-password", r.resetPassword)
		usergrp.GET("refresh-token", r.refreshAuthToken)
	}
}

func (r *authroutes) login(
	ctx *gin.Context,
) {
	var function = "login"
	var req requests.Login

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

	user, token, err := r.sauth.Login(req, ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	res := token.ToAuthResponse(*user)

	utresponse.GeneralSuccessAuth(
		function,
		ctx,
		res,
		token,
	)
}

func (r *authroutes) register(
	ctx *gin.Context,
) {
	var function = "register"
	var req requests.Register

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

	user, token, err := r.sauth.Register(req, ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	res := token.ToAuthResponse(*user)

	utresponse.GeneralSuccessAuth(
		function,
		ctx,
		res,
		token,
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

	user, token, err := r.sauth.RefreshAuthToken(refreshToken, ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	res := token.ToAuthResponse(*user)

	utresponse.GeneralSuccessAuth(
		function,
		ctx,
		res,
		token,
	)
}
