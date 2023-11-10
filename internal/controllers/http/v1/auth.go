package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	authservice "project-skbackend/internal/services/auth"
	"project-skbackend/packages/utils"

	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	cfg     *configs.Config
	authsvc authservice.IAuthService
}

func newAuthRoutes(
	h *gin.RouterGroup,
	cfg *configs.Config,
	authsvc authservice.IAuthService,
) {
	r := &authRoutes{
		cfg:     cfg,
		authsvc: authsvc,
	}

	useg := h.Group("auth")
	{
		useg.POST("login", r.login)
		useg.POST("register", r.register)

		verifyGroup := useg.Group("verify")
		{
			verifyGroup.POST("", r.verifyToken)
			verifyGroup.POST("send", r.sendVerifyEmail)
		}

		useg.POST("forgot-password", r.forgotPassword)
		useg.POST("reset-password", r.resetPassword)

		useg.GET("refresh-token", r.refreshAuthToken)
	}
}

func (r *authRoutes) login(
	ctx *gin.Context,
) {
	var req requests.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utils.ValidationResponse(err)

		utils.GeneralInputRequiredError(
			"login",
			ctx,
			ve,
		)
		return
	}

	user, token, err := r.authsvc.Login(req)
	if err != nil {
		utils.GeneralInternalServerError(
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

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Login Successful",
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
		ve := utils.ValidationResponse(err)

		utils.GeneralInputRequiredError(
			"register",
			ctx,
			ve,
		)
		return
	}

	user, token, err := r.authsvc.Register(req)
	if err != nil {
		utils.GeneralInternalServerError(
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

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
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
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  utils.ErrUserNotFound,
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  utils.ErrUserIDNotFound,
		})
		return
	}

	token := utils.GenerateRandomToken()
	err := r.authsvc.SendVerificationEmail(loggedInUser.ID, token)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error sending verification email",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
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
		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ErrorRes{
			Message: "Invalid request",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	err = r.authsvc.VerifyToken(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ErrorRes{
			Message: "Cannot verify token",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
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
		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ErrorRes{
			Message: "Invalid request",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	err = r.authsvc.ForgotPassword(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
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
		ve := utils.ValidationResponse(err)

		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ErrorRes{
			Message: "Invalid request",
			Debug:   err,
			Errors:  ve,
		})
		return
	}

	err = r.authsvc.ResetPassword(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, utils.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Your password has been successfully changed",
		Data:    nil,
	})
}

func (r *authRoutes) refreshAuthToken(
	ctx *gin.Context,
) {
	refreshToken := ctx.Request.Header.Get("Refresh-Token")

	if refreshToken == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.ErrorRes{
			Message: "Something went wrong",
			Debug:   nil,
			Errors:  "No Refresh token detected",
		})
		return
	}

	user, token, err := r.authsvc.RefreshAuthToken(refreshToken)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, utils.ErrorRes{
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

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Refresh Token Successful",
		Data:    res,
		Header:  *token,
	})
}
