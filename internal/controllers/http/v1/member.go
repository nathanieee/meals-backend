package controllers

import (
	"fmt"
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
	var function = "create member"
	var entity = "member"
	var req requests.CreateMember

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	// * get the user image
	file := req.User.UserImage
	if file != nil {
		// * check if the file is an image
		err := file.IsImage()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	meres, err := r.smember.Create(req)
	if err != nil {
		fmt.Println(err, "create member")
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			fmt.Println(err)
			utresponse.GeneralDuplicate(
				"email",
				ctx,
				err,
			)
		} else {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
		}
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		meres,
	)
}

func (r *memberroutes) getMembers(ctx *gin.Context) {
	var entity = "members"
	paginationReq := utrequest.GeneratePaginationFromRequest(ctx)

	members, err := r.smember.FindAll(paginationReq)
	if err != nil {
		utresponse.GeneralNotFound(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetching(
		entity,
		ctx,
		members,
	)
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
