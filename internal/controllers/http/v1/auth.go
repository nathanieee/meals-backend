package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
)

type (
	authRoutes struct {
		cfg   *configs.Config
		sauth authservice.IAuthService
	}
)

func newAuthRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
) {
	r := &authRoutes{
		cfg:   cfg,
		sauth: sauth,
	}

	usergrp := rg.Group("auth")
	{
		usergrp.POST("login", r.login)
		usergrp.POST("register", r.register)

		verifgrp := usergrp.Group("verify")
		{
			verifgrp.POST("", r.verifyToken)
			verifgrp.POST("send", r.sendVerifyEmail)
		}

		usergrp.POST("forgot-password", r.forgotPassword)
		usergrp.POST("reset-password", r.resetPassword)
		usergrp.GET("refresh-token", r.refreshAuthToken)
	}
}

func (r *authRoutes) login(
	ctx *gin.Context,
) {
	var req requests.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInputRequiredError(
			"login",
			ctx,
			ve,
		)
		return
	}

	user, token, err := r.sauth.Login(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			"login",
			ctx,
			err.Error(),
		)
		return
	}

	res := responses.AuthResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Role:               user.Role,
		ConfirmationSentAt: user.ConfirmationSentAt,
		ConfirmedAt:        user.ConfirmedAt,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		Token:              token.AuthToken,
		Expires:            token.AuthTokenExpires,
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "login successful",
		Data:    res,
		Header:  *token,
	})
}

func (r *authRoutes) register(
	ctx *gin.Context,
) {
	var req requests.RegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInputRequiredError(
			"register",
			ctx,
			ve,
		)
		return
	}

	user, token, err := r.sauth.Register(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			"register",
			ctx,
			err.Error(),
		)
		return
	}

	res := responses.AuthResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Role:               user.Role,
		ConfirmationSentAt: user.ConfirmationSentAt,
		ConfirmedAt:        user.ConfirmedAt,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		Token:              token.AuthToken,
		Expires:            token.AuthTokenExpires,
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Register Successful",
		Data:    res,
		Header:  *token,
	})
}

func (r *authRoutes) sendVerifyEmail(
	ctx *gin.Context,
) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  utresponse.ErrUserNotFound,
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  utresponse.ErrUserIDNotFound,
		})
		return
	}

	token := uttoken.GenerateRandomToken()
	err := r.sauth.SendVerificationEmail(loggedInUser.ID, token)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error sending verification email",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Send Verification Email Successful",
		Data:    nil,
	})
}

func (r *authRoutes) verifyToken(
	ctx *gin.Context,
) {
	var req requests.VerifyTokenRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Message: "Invalid request",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	err = r.sauth.VerifyToken(req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Message: "Cannot verify token",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Verification successful",
		Data:    nil,
	})
}

func (r *authRoutes) forgotPassword(
	ctx *gin.Context,
) {
	var req requests.ForgotPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Message: "Invalid request",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	err = r.sauth.ForgotPassword(req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Successfully requested to forgot password",
		Data:    nil,
	})
}

func (r *authRoutes) resetPassword(
	ctx *gin.Context,
) {
	var req requests.ResetPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Message: "Invalid request",
			Debug:   err,
			Errors:  ve,
		})
		return
	}

	err = r.sauth.ResetPassword(req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Your password has been successfully changed",
		Data:    nil,
	})
}

func (r *authRoutes) refreshAuthToken(
	ctx *gin.Context,
) {
	refreshToken := ctx.Request.Header.Get("Refresh-Token")

	if refreshToken == "" {
		utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
			Message: "Something went wrong",
			Debug:   nil,
			Errors:  "No Refresh token detected",
		})
		return
	}

	user, token, err := r.sauth.RefreshAuthToken(refreshToken)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	res := responses.AuthResponse{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		ConfirmedAt: user.ConfirmedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Token:       token.AuthToken,
		Expires:     token.AuthTokenExpires,
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Refresh Token Successful",
		Data:    res,
		Header:  *token,
	})
}
