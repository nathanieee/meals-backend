package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/utstring"

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

	user, token, err := r.sauth.Login(req, ctx)
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
		Token:              token.AccessToken,
		Expires:            token.AccessTokenExpires,
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
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

	user, token, err := r.sauth.Register(req, ctx)
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
		Token:              token.AccessToken,
		Expires:            token.AccessTokenExpires,
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "register successful",
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
			Status:  consttypes.RST_FAIL,
			Message: "Error getting user",
			Data: utresponse.ErrorData{
				Debug:  nil,
				Errors: utresponse.ErrUserNotFound,
			},
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error getting user",
			Data: utresponse.ErrorData{
				Debug:  nil,
				Errors: utresponse.ErrUserIDNotFound,
			},
		})
		return
	}

	token, err := utstring.GenerateRandomToken(r.cfg.VerifyTokenLength)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusInternalServerError, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error generating token",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	err = r.sauth.SendVerificationEmail(loggedInUser.ID, token)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error sending verification email",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "send verification email successful",
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
			Status:  consttypes.RST_FAIL,
			Message: "Invalid request",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	err = r.sauth.VerifyToken(req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Status:  consttypes.RST_FAIL,
			Message: "Cannot verify token",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "verification successful",
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
			Status:  consttypes.RST_FAIL,
			Message: "Invalid request",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	err = r.sauth.ForgotPassword(req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Something went wrong",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
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

		utresponse.GeneralInvalidRequest(
			"reset password",
			ctx,
			ve,
			&err,
		)
		return
	}

	err = r.sauth.ResetPassword(req)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusBadRequest, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Something went wrong",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "Your password has been successfully changed",
		Data:    nil,
	})
}

func (r *authRoutes) refreshAuthToken(
	ctx *gin.Context,
) {
	refreshToken, err := ctx.Cookie("refresh_token")

	if refreshToken == "" || err != nil {
		utresponse.ErrorResponse(ctx, http.StatusUnauthorized, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Something went wrong",
			Data: utresponse.ErrorData{
				Debug:  nil,
				Errors: "No Refresh token detected",
			},
		})
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

	res := responses.AuthResponse{
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
