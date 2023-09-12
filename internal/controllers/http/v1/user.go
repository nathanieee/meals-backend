package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/services"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userRoutes struct {
	us  services.IUserService
	ms  services.IMailService
	cfg *configs.Config
}

func newUserRoutes(h *gin.RouterGroup, db *gorm.DB, cfg *configs.Config, us services.IUserService, ms services.IMailService) {
	r := &userRoutes{us: us, cfg: cfg, ms: ms}

	adminGroup := h.Group("users")
	{
		adminGroup.GET("", r.getUser)
		adminGroup.POST("", r.createUser)
	}

	userGroup := h.Group("users")
	{
		userGroup.GET("/me", r.getCurrentUser)
		userGroup.DELETE("/delete", r.deleteUser)
	}
}

func (r *userRoutes) createUser(ctx *gin.Context) {
	var req requests.CreateUserRequest

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

	user, err := r.us.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			utils.ErrorResponse(ctx, http.StatusConflict, utils.ErrorRes{
				Message: "Duplicate email",
				Debug:   err,
				Errors:  err.Error(),
			})
		} else {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, utils.ErrorRes{
				Message: "Something went wrong",
				Debug:   err,
				Errors:  err.Error(),
			})
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Success Creating new user",
		Data:    user,
	})
}

func (r *userRoutes) getUser(ctx *gin.Context) {
	paginationReq := utils.GeneratePaginationFromRequest(ctx, models.User{})

	users, err := r.us.GetUsers(paginationReq)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "User not found",
			Debug:   nil,
			Errors:  nil,
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Success Get Users",
		Data:    users,
	})
}

func (r *userRoutes) getCurrentUser(ctx *gin.Context) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "User not found",
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "Unable to assert User ID",
		})
		return
	}

	user, err := r.us.GetUser(loggedInUser.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "User not found",
			Debug:   nil,
			Errors:  "User not found",
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Success Get User",
		Data:    user,
	})
}

func (r *userRoutes) deleteUser(ctx *gin.Context) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "User not found",
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.UserResponse)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Error getting user",
			Debug:   nil,
			Errors:  "Unable to assert User ID",
		})
		return
	}

	if loggedInUser.Role == consttypes.UR_ADMINISTRATOR {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Something went wrong",
			Debug:   nil,
			Errors:  "Admin role can't be deleted",
		})
		return
	}

	err := r.us.DeleteUser(loggedInUser.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, utils.ErrorRes{
			Message: "Something went wrong",
			Debug:   err,
			Errors:  "Something went wrong while deleting user",
		})
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, utils.SuccessRes{
		Message: "Success Delete User",
		Data:    nil,
	})
}
