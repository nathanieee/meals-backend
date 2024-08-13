package controllers

import (
	"errors"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/internal/services/patronservice"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	patronroutes struct {
		cfg     *configs.Config
		spatron patronservice.IPatronService
		sauth   authservice.IAuthService
		sfile   fileservice.IFileService
	}
)

func newPatronRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
	spatron patronservice.IPatronService,
	sfile fileservice.IFileService,
) {
	r := &patronroutes{
		cfg:     cfg,
		sauth:   sauth,
		spatron: spatron,
		sfile:   sfile,
	}

	gpatronspub := rg.Group("patrons")
	{
		gpatronspub.POST("register", r.patronRegister)
	}
}

func (r *patronroutes) patronRegister(ctx *gin.Context) {
	var (
		function = "patron register"
		entity   = "patron"
		req      requests.CreatePatron
		err      error
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

	respatron, err := r.spatron.Create(req)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgerr.SQLState()) {
				utresponse.GeneralDuplicate(
					pgerr.TableName,
					ctx,
					pgerr,
				)
				return
			}
		} else {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
		}
		return
	}

	// * define the image request
	reqimg := req.User.CreateImage
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

		err = r.sfile.UploadProfilePicture(respatron.User.ID, multipart)
		if err != nil {
			utresponse.GeneralInternalServerError(
				function,
				ctx,
				err,
			)
			return
		}
	}

	resuser, thead, err := r.sauth.Signin(*req.ToSignin(), ctx)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	err = r.sauth.SendVerificationEmail(resuser.ID)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	resauth := thead.ToAuthResponse(*resuser)
	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resauth,
	)
}
