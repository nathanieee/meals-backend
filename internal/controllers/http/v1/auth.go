package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
)

type (
	authroutes struct {
		cfg   *configs.Config
		sauth authservice.IAuthService
		suser userservice.IUserService
	}
)

func newAuthRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	sauth authservice.IAuthService,
	suser userservice.IUserService,
) {
	r := &authroutes{
		cfg:   cfg,
		sauth: sauth,
		suser: suser,
	}

	h := rg.Group("auth")
	{
		h.POST("signin", r.signin)

		gverif := h.Group("verify").Use(middlewares.JWTAuthMiddleware(cfg,
			uint(consttypes.UR_ADMIN),
			uint(consttypes.UR_CAREGIVER),
			uint(consttypes.UR_MEMBER),
			uint(consttypes.UR_ORGANIZATION),
			uint(consttypes.UR_PARTNER),
			uint(consttypes.UR_PATRON),
			uint(consttypes.UR_USER),
		))
		{
			gverif.POST("", r.verifyToken)
			gverif.POST("send", r.sendVerifyEmail)
		}

		h.POST("forgot-password", r.forgotPassword)
		h.POST("reset-password", r.resetPassword)
		h.GET("refresh-token", r.refreshAuthToken)
	}
}

func (r *authroutes) signin(
	ctx *gin.Context,
) {
	var (
		function = "signin"
		req      requests.Signin
	)

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

	resuser, thead, err := r.sauth.Signin(req, ctx)
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

func (r *authroutes) forgotPassword(
	ctx *gin.Context,
) {
	var (
		function = "forgot password"
		req      requests.ForgotPassword
	)

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

	err = r.sauth.ForgotPassword(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccess(
		function,
		ctx,
		nil,
	)
}

func (r *authroutes) resetPassword(
	ctx *gin.Context,
) {
	var (
		function = "reset password"
		req      requests.ResetPassword
	)

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

	err = r.sauth.ResetPassword(req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccess(
		function,
		ctx,
		nil,
	)
}

func (r *authroutes) refreshAuthToken(
	ctx *gin.Context,
) {
	var (
		function = "refresh token"
	)

	trefresh, err := ctx.Cookie("refresh-token")
	if trefresh == "" || err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			utresponse.ErrTokenNotFound,
		)
		return
	}

	resuser, thead, err := r.sauth.RefreshAuthToken(trefresh, ctx)
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

func (r *authroutes) verifyToken(ctx *gin.Context) {
	var (
		function = "verify token"
		req      requests.VerifyToken
	)

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

	resuser, theader, err := r.sauth.VerifyToken(req, ctx)
	if err != nil {
		utresponse.GeneralInvalidRequest(
			function,
			ctx,
			nil,
			err,
		)
		return
	}

	resauth := theader.ToAuthResponse(*resuser)
	utresponse.GeneralSuccessAuth(
		function,
		ctx,
		resauth,
		theader,
	)
}

func (r *authroutes) sendVerifyEmail(ctx *gin.Context) {
	var (
		function = "send verify email"
	)

	user, err := uttoken.GetUser(ctx)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	err = r.sauth.SendVerificationEmail(user.ID)
	if err != nil {
		var (
			entity = "user"
		)
		utresponse.GeneralNotFound(
			entity,
			ctx,
			utresponse.ErrUserNotFound,
		)
		return
	}

	utresponse.GeneralSuccess(
		function,
		ctx,
		nil,
	)
}
