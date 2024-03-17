package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
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
	cfg *configs.Config,
	smember memberservice.IMemberService,
) {
	r := &memberroutes{
		cfg:     cfg,
		smember: smember,
	}

	gadmn := rg.Group("members")
	gadmn.Use(middlewares.JWTAuthMiddleware(
		cfg,
		uint(consttypes.UR_ADMIN),
	))
	{
		gadmn.POST("", r.createMember)
		gadmn.GET("", r.getMembers)
		gadmn.GET("raw", r.getMembersRaw)
		gadmn.PUT("/:uuid", r.updateMember)
		gadmn.DELETE("/:uuid", r.deleteMember)
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
	file := req.User.CreateImage
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

	resmemb, err := r.smember.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
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
		resmemb,
	)
}

func (r *memberroutes) getMembers(ctx *gin.Context) {
	var entity = "members"
	var reqpage = utrequest.GeneratePaginationFromRequest(ctx)

	members, err := r.smember.FindAll(reqpage)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		members,
	)
}

func (r *memberroutes) getMembersRaw(ctx *gin.Context) {
	var entity = "members"

	resmemb, err := r.smember.Read()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		entity,
		ctx,
		resmemb,
	)
}

func (r *memberroutes) updateMember(ctx *gin.Context) {
	var function = "update member"
	var entity = "member"
	var req requests.UpdateMember

	err := ctx.ShouldBind(&req)
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

	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smember.FindByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	resmemb, err := r.smember.Update(uuid, req)
	if err != nil {
		utresponse.GeneralFailedUpdate(
			entity,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resmemb,
	)
}

func (r *memberroutes) deleteMember(ctx *gin.Context) {
	var function = "delete member"
	var entity = "member"

	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

	_, err = r.smember.FindByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utresponse.GeneralNotFound(
				entity,
				ctx,
				err,
			)
			return
		}

		utresponse.GeneralInternalServerError(
			entity,
			ctx,
			err,
		)
		return
	}

	err = r.smember.Delete(uuid)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
	}

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}
