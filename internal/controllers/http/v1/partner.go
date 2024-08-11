package controllers

import (
	"errors"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/internal/services/partnerservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	partnerroutes struct {
		cfg      *configs.Config
		sauth    authservice.IAuthService
		spartner partnerservice.IPartnerService
		sfile    fileservice.IFileService
	}
)

func newPartnerRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
	spartner partnerservice.IPartnerService,
	sfile fileservice.IFileService,
) {
	r := &partnerroutes{
		cfg:      cfg,
		sauth:    sauth,
		spartner: spartner,
		sfile:    sfile,
	}

	gpartnerspub := rg.Group("partners")
	{
		gpartnerspub.POST("register", r.partnerRegister)
	}

	gpartnerspvt := rg.Group("partners")
	gpartnerspvt.Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_PARTNER))
	{

	}
}

func (r *partnerroutes) partnerRegister(ctx *gin.Context) {
	var (
		function = "partner register"
		req      requests.CreatePartner
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

	respartner, err := r.spartner.Create(req)
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

		err = r.sfile.UploadProfilePicture(respartner.User.ID, multipart)
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
	utresponse.GeneralSuccessAuth(
		function,
		ctx,
		resauth,
		thead,
	)
}
