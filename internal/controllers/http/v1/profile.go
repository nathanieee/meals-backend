package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/baseroleservice"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	profileroutes struct {
		cfg   *configs.Config
		suser userservice.IUserService
		smemb memberservice.IMemberService
		sfile fileservice.IFileService
		sbase baseroleservice.IBaseRoleService
	}
)

func newProfileRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	suser userservice.IUserService,
	smemb memberservice.IMemberService,
	sfile fileservice.IFileService,
	sbase baseroleservice.IBaseRoleService,
) {
	r := &profileroutes{
		cfg:   cfg,
		suser: suser,
		smemb: smemb,
		sfile: sfile,
		sbase: sbase,
	}

	gprofilepvt := rg.Group("profiles")
	gprofilepvt.Use(middlewares.JWTAuthMiddleware(
		cfg,
		consttypes.UR_MEMBER,
		consttypes.UR_CAREGIVER,
		consttypes.UR_ADMIN,
		consttypes.UR_ORGANIZATION,
		consttypes.UR_PARTNER,
		consttypes.UR_PATRON,
	))
	{
		gprofilemem := gprofilepvt.Group("members")
		gprofilemem.Use(middlewares.JWTAuthMiddleware(
			cfg,
			consttypes.UR_MEMBER,
		))
		{
			gprofilemem.PATCH("own", r.updateOwnMemberProfile)
		}

		gprofilepvt.GET("me", r.getOwnProfile)
	}
}

func (r *profileroutes) getOwnProfile(ctx *gin.Context) {
	var (
		function = "get own profile"
	)

	userres, err := uttoken.GetUser(ctx)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	roleres, err := r.suser.GetRoleDataByUserID(userres.ID)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	if roleres == nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			consttypes.ErrUserInvalidRole,
		)
		return
	}

	utresponse.GeneralSuccessFetch(
		function,
		ctx,
		roleres.Data,
	)
}

func (r *profileroutes) updateOwnMemberProfile(ctx *gin.Context) {
	var (
		function = "update member profile"
		entity   = "member profile"
		req      *requests.UpdateMember
	)

	if err := ctx.ShouldBind(&req); err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			ve,
			err,
		)
		return
	}

	// * get the current logged in user
	userres, err := uttoken.GetUser(ctx)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	// * get its role data by its user id
	roleres, err := r.suser.GetRoleDataByUserID(userres.ID)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	// * validate if the role data is empty
	if roleres == nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			consttypes.ErrUserInvalidRole,
		)
		return
	}

	modmem, err := r.sbase.GetMemberByBaseRole(*roleres)
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

	resmem, err := r.smemb.Update(modmem.ID, *req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	// * define the image request
	reqimg := req.User.UpdateImage
	// * if the image request is not empty
	// * validate and upload the image
	if reqimg != nil {
		if err := reqimg.Validate(); err != nil {
			utresponse.GeneralInvalidRequest(
				function,
				ctx,
				nil,
				err,
			)
			return
		}

		multipart, err := reqimg.GetMultipartFile()
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}

		err = r.sfile.UploadProfilePicture(resmem.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	resmem, err = r.smemb.GetByID(resmem.ID)
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

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		resmem,
	)
}
