package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	userRoutes struct {
		cfg   *configs.Config
		suser userservice.IUserService
		smail mailservice.IMailService
	}
)

func newUserRoutes(
	rg *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	suser userservice.IUserService,
	smail mailservice.IMailService,
) {
	r := &userRoutes{
		cfg:   cfg,
		suser: suser,
		smail: smail,
	}

	admingrp := rg.Group("users")
	{
		admingrp.GET("", r.getUser)
		admingrp.POST("", r.createUser)
	}

	usergrp := rg.Group("users")
	usergrp.Use(middlewares.JWTAuthMiddleware(
		cfg,
		uint(consttypes.UR_USER),
	))
	{
		usergrp.GET("/me", r.getCurrentUser)
		usergrp.DELETE("/delete", r.deleteUser)
	}
}

func (r *userRoutes) createUser(ctx *gin.Context) {
	var req requests.CreateUserRequest

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

	ures, err := r.suser.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			utresponse.ErrorResponse(ctx, http.StatusConflict, utresponse.ErrorRes{
				Message: "Duplicate email",
				Debug:   err,
				Errors:  err.Error(),
			})
		} else {
			utresponse.ErrorResponse(ctx, http.StatusInternalServerError, utresponse.ErrorRes{
				Message: "Something went wrong",
				Debug:   err,
				Errors:  err.Error(),
			})
		}
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Creating new user",
		Data:    ures,
	})
}

func (r *userRoutes) getUser(ctx *gin.Context) {
	paginationReq := utrequest.GeneratePaginationFromRequest(ctx)

	users, err := r.suser.FindAll(paginationReq)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "users not found",
			Debug:   nil,
			Errors:  nil,
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Get Users",
		Data:    users,
	})
}

func (r *userRoutes) getCurrentUser(ctx *gin.Context) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "User not found",
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "Unable to assert User ID",
		})
		return
	}

	user, err := r.suser.FindByID(loggedInUser.ID)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "User not found",
			Debug:   nil,
			Errors:  "User not found",
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Get User",
		Data:    user,
	})
}

func (r *userRoutes) deleteUser(ctx *gin.Context) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "User not found",
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "Unable to assert User ID",
		})
		return
	}

	if loggedInUser.Role == consttypes.UR_ADMIN {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Something went wrong",
			Debug:   nil,
			Errors:  "Admin role can't be deleted",
		})
		return
	}

	err := r.suser.Delete(loggedInUser.ID)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  "Something went wrong while deleting user",
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Delete User",
		Data:    nil,
	})
}
