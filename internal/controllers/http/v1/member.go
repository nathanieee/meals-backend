package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	memberRoutes struct {
		cfg     *configs.Config
		smember memberservice.IMemberService
	}
)

func newMemberRoutes(
	rg *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	smember memberservice.IMemberService,
) {
	r := &memberRoutes{
		cfg:     cfg,
		smember: smember,
	}

	admgrp := rg.Group("members")
	{
		admgrp.POST("", r.createMember)
		admgrp.GET("", r.getMembers)
	}
}

func (r *memberRoutes) createMember(ctx *gin.Context) {
	var req requests.CreateMember

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)

		utresponse.GeneralInvalidRequest(
			"create member",
			ctx,
			ve,
			&err,
		)
		return
	}

	meres, err := r.smember.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			utresponse.ErrorResponse(ctx, http.StatusConflict, utresponse.ErrorRes{
				Status:  consttypes.RST_ERROR,
				Message: "duplicate email",
				Data: utresponse.ErrorData{
					Debug:  err,
					Errors: err.Error(),
				},
			})
		} else {
			utresponse.GeneralInternalServerError(
				"something went wrong",
				ctx,
				err,
			)
		}
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "success creating new member",
		Data:    meres,
	})
}

func (r *memberRoutes) getMembers(ctx *gin.Context) {
	paginationReq := utrequest.GeneratePaginationFromRequest(ctx)

	members, err := r.smember.FindAll(paginationReq)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "members not found",
			Data: utresponse.ErrorData{
				Debug:  err,
				Errors: err.Error(),
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "success get members",
		Data:    members,
	})
}
