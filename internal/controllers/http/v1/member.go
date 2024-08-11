package controllers

import (
	"errors"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/orderservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	memberroutes struct {
		cfg     *configs.Config
		smember memberservice.IMemberService
		scart   cartservice.ICartService
		suser   userservice.IUserService
		sauth   authservice.IAuthService
		sorder  orderservice.IOrderService
		sfile   fileservice.IFileService
	}
)

func newMemberRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	smember memberservice.IMemberService,
	scart cartservice.ICartService,
	suser userservice.IUserService,
	sauth authservice.IAuthService,
	sorder orderservice.IOrderService,
	sfile fileservice.IFileService,
) {
	r := &memberroutes{
		cfg:     cfg,
		smember: smember,
		scart:   scart,
		suser:   suser,
		sauth:   sauth,
		sorder:  sorder,
		sfile:   sfile,
	}

	gmemberspub := rg.Group("members")
	{
		gmemberspub.POST("register", r.memberRegister)
	}

	gmemberspvt := rg.Group("members")
	gmemberspvt.Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_MEMBER))
	{
		gcart := gmemberspvt.Group("carts")
		{
			gcart.POST("", r.memberCreateCart)
		}

		gorder := gmemberspvt.Group("orders")
		{
			gorder.POST("", r.memberCreateOrder)
		}
	}
}

func (r *memberroutes) memberRegister(ctx *gin.Context) {
	var (
		function = "member register"
		entity   = "member"
		req      requests.CreateMember
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

	member, err := r.smember.Create(req)
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

		err = r.sfile.UploadProfilePicture(member.User.ID, multipart)
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
		entity,
		ctx,
		resauth,
		thead,
	)
}

func (r *memberroutes) memberCreateCart(ctx *gin.Context) {
	var (
		function = "create cart"
		entity   = "cart"
		req      *requests.CreateCart
		err      error
	)

	userres, err := uttoken.GetUser(ctx)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

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

	rescart, err := r.scart.Create(*req, *roleres)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		rescart,
	)
}

func (r *memberroutes) memberCreateOrder(ctx *gin.Context) {
	var (
		function = "create order"
		entity   = "order"
		req      *requests.CreateOrder
		err      error
	)

	userres, err := uttoken.GetUser(ctx)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

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

	resorder, err := r.sorder.Create(*req, userres.ID)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessCreate(
		entity,
		ctx,
		resorder,
	)
}
