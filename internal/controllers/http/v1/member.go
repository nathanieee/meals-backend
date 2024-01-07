package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	mmbrservice "project-skbackend/internal/services/member"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type memberRoutes struct {
	cfg     *configs.Config
	membsvc mmbrservice.IMemberService
}

func newMemberRoutes(
	rg *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	membsvc mmbrservice.IMemberService,
) {
	r := &memberRoutes{
		cfg:     cfg,
		membsvc: membsvc,
	}

	admgrp := rg.Group("members")
	{
		admgrp.POST("", r.createMember)
		admgrp.GET("", r.getMembers)
	}
}

func (r *memberRoutes) createMember(ctx *gin.Context) {
	var req requests.CreateMemberRequest

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

	meres, err := r.membsvc.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			utresponse.ErrorResponse(ctx, http.StatusConflict, utresponse.ErrorRes{
				Message: "duplicate email",
				Debug:   err,
				Errors:  err.Error(),
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
		Message: "success creating new member",
		Data:    meres,
	})
}

func (r *memberRoutes) getMembers(ctx *gin.Context) {
	paginationReq := utrequest.GeneratePaginationFromRequest(ctx)

	members, err := r.membsvc.FindAll(paginationReq)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Message: "members not found",
			Debug:   err,
			Errors:  err.Error(),
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "success get members",
		Data:    members,
	})
}
