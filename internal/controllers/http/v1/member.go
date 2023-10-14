package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type memberRoutes struct {
	mes services.IMemberService
	cfg *configs.Config
}

func newMemberRoutes(
	h *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	mes services.IMemberService,
) {
	r := &memberRoutes{
		mes: mes,
		cfg: cfg,
	}

	adminGroup := h.Group("members")
	adminGroup.Use(middlewares.JWTAuthMiddleware(
		cfg,
		uint(consttypes.UR_USER),
	))
	{
		adminGroup.POST("", r.createMember)
	}
}

func (r *memberRoutes) createMember(ctx *gin.Context) {
	var req requests.CreateMemberRequest

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

	meres, err := r.mes.Create(req)
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
		Data:    meres,
	})
}
