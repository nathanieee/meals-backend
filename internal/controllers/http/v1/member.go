package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	memberroutes struct {
		cfg     *configs.Config
		smember memberservice.IMemberService
		scart   cartservice.ICartService
		suser   userservice.IUserService
		sauth   authservice.IAuthService
	}
)

func newMemberRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	smember memberservice.IMemberService,
	scart cartservice.ICartService,
	suser userservice.IUserService,
	sauth authservice.IAuthService,
) {
	r := &memberroutes{
		cfg:     cfg,
		smember: smember,
		scart:   scart,
		suser:   suser,
		sauth:   sauth,
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
	}
}

func (r *memberroutes) memberRegister(ctx *gin.Context) {
	var (
		function = "member register"
		entity   = "member"
		req      requests.CreateMember
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	_, err = r.smember.Create(req)
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

	resuser, thead, err := r.sauth.Signin(*req.ToSignin(), ctx)
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
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
