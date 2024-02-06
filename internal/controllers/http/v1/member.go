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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	memberroutes struct {
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
	r := &memberroutes{
		cfg:     cfg,
		smember: smember,
	}

	admgrp := rg.Group("members")
	{
		admgrp.POST("", r.createMember)
		admgrp.GET("", r.getMembers)
		admgrp.PUT("/:uuid", r.updateMember)
	}
}

func (r *memberroutes) createMember(ctx *gin.Context) {
	var req requests.CreateMember

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			"create member failed",
			ctx,
			ve,
			err,
		)
		return
	}

	// * get the user image
	file := req.User.Image
	if file != nil {
		// * check if the file is an image
		err := file.IsImage()
		if err != nil {
			utresponse.GeneralInvalidRequest(
				"file validation failed",
				ctx,
				nil,
				err,
			)
			return
		}
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
		Message: "success creating a new member",
		Data:    meres,
	})
}

func (r *memberroutes) getMembers(ctx *gin.Context) {
	paginationReq := utrequest.GeneratePaginationFromRequest(ctx)

	members, err := r.smember.FindAll(paginationReq)
	if err != nil {
		utresponse.GeneralNotFound(
			"members",
			ctx,
			err,
		)
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "success get members",
		Data:    members,
	})
}

func (r *memberroutes) updateMember(ctx *gin.Context) {
	var req requests.UpdateMember

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			"member",
			ctx,
			ve,
			err,
		)
		return
	}

	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			utresponse.ErrConvertFailed.Error(),
			ctx,
			err,
		)
		return
	}

	_, err = r.smember.FindByID(uuid)
	if err != nil {
		utresponse.GeneralNotFound(
			"member",
			ctx,
			err,
		)
		return
	}

	mres, err := r.smember.Update(uuid, req)
	if err != nil {
		utresponse.GeneralFailedUpdate(
			"member",
			ctx,
			err,
		)
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: "success update member",
		Data:    mres,
	})
}
