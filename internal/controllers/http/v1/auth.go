package controllers

import (
	"fmt"
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/packages/consttypes"
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
	var req requests.Login

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInvalidRequest(
			"login",
			ctx,
			ve,
			err,
		)
		return
	}

	user, token, err := r.sauth.Login(req, ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			"login",
			ctx,
			err,
		)
		return
	}

	res := token.ToAuthResponse(*user)

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "login successful",
		Data:    res,
		Header:  *token,
	})
}

func (r *authroutes) register(
	ctx *gin.Context,
) {
	var req requests.Register

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInvalidRequest(
			"register",
			ctx,
			ve,
			err,
		)
		return
	}

	user, token, err := r.sauth.Register(req, ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			"register",
			ctx,
			err,
		)
		return
	}

	res := token.ToAuthResponse(*user)

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "register successful",
		Data:    res,
		Header:  *token,
	})
}

func (r *authroutes) forgotPassword(
	ctx *gin.Context,
) {
	var req requests.ForgotPassword

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInvalidRequest(
			"forgot password",
			ctx,
			ve,
			err,
		)
		return
	}

	err = r.sauth.ForgotPassword(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			"forgot password",
			ctx,
			err,
		)
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "Successfully requested to forgot password",
		Data:    nil,
	})
}

func (r *authroutes) resetPassword(
	ctx *gin.Context,
) {
	var req requests.ResetPassword

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInvalidRequest(
			"reset password",
			ctx,
			ve,
			err,
		)
		return
	}

	err = r.sauth.ResetPassword(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			"reset password",
			ctx,
			err,
		)
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "Your password has been successfully changed",
		Data:    nil,
	})
}

func (r *authroutes) refreshAuthToken(
	ctx *gin.Context,
) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if refreshToken == "" || err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			fmt.Errorf("refresh token not found"),
		)
		return
	}

	user, token, err := r.sauth.RefreshAuthToken(refreshToken, ctx)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Something went wrong",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	res := responses.Auth{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		ConfirmedAt: user.ConfirmedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Token:       token.AccessToken,
		Expires:     token.AccessTokenExpires,
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "Refresh Token Successful",
		Data:    res,
		Header:  *token,
	})
}
