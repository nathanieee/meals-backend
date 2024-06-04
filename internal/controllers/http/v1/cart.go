package controllers

import (
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utresponse"
	"project-skbackend/packages/utils/utrole"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	cartroutes struct {
		cfg   *configs.Config
		scart cartservice.ICartService
		suser userservice.IUserService
		smmbr memberservice.IMemberService
	}
)

func newCartRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	scart cartservice.ICartService,
	suser userservice.IUserService,
	smmbr memberservice.IMemberService,
) {
	r := &cartroutes{
		cfg:   cfg,
		scart: scart,
		suser: suser,
		smmbr: smmbr,
	}

	gcart := rg.Group("carts")
	{

		gmembcare := gcart.Group("").Use(middlewares.JWTAuthMiddleware(cfg, consttypes.UR_MEMBER, consttypes.UR_CAREGIVER))
		{
			gmembcare.POST("", r.createCart)
			gmembcare.DELETE(":muuid", r.deleteCart)
			gmembcare.GET("raw", r.getCartsRaw)
		}

	}
}

func (r *cartroutes) createCart(ctx *gin.Context) {
	var (
		function = "create cart"
		entity   = "cart"
		req      *requests.CreateCart
		err      error

		rid  uuid.UUID
		role consttypes.UserRole
		ok   bool
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

	rid, role, ok = utrole.CartRoleCheck(*roleres)
	if !ok {
		utresponse.GeneralUnauthorized(
			ctx,
			consttypes.ErrUserInvalidRole,
		)
		return
	}

	if role == consttypes.UR_CAREGIVER {
		m, err := r.smmbr.FindByCaregiverID(rid)
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

		rid = m.ID
		role = m.User.Role
	}

	req, err = req.New(rid, role)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	rescart, err := r.scart.Create(*req)
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

func (r *cartroutes) getCartsRaw(ctx *gin.Context) {
	var (
		function = "get carts raw"
		entity   = "carts"

		rid  uuid.UUID
		role consttypes.UserRole
		ok   bool
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

	rid, role, ok = utrole.CartRoleCheck(*roleres)
	if !ok {
		utresponse.GeneralUnauthorized(
			ctx,
			consttypes.ErrUserInvalidRole,
		)
		return
	}

	if role == consttypes.UR_CAREGIVER {
		m, err := r.smmbr.FindByCaregiverID(rid)
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

		rid = m.ID
		role = m.User.Role
	}

	carts, err := r.scart.ReadWithReference(rid, role)
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
		carts,
	)
}

// * unused function
func (r *cartroutes) updateCart(ctx *gin.Context) {
	var (
		function = "update cart"
		entity   = "cart"
		req      *requests.UpdateCart
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

	roleres, err := r.suser.GetRoleDataByUserID(userres.ID)
	if err != nil {
		utresponse.GeneralUnauthorized(
			ctx,
			err,
		)
		return
	}

	if roleres != nil {
		_, _, ok := utrole.CartRoleCheck(*roleres)
		if !ok {
			utresponse.GeneralUnauthorized(
				ctx,
				consttypes.ErrUserInvalidRole,
			)
			return
		}
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

	rescart, err := r.scart.Update(*req)
	if err != nil {
		utresponse.GeneralInternalServerError(
			function,
			ctx,
			err,
		)
		return
	}

	utresponse.GeneralSuccessUpdate(
		entity,
		ctx,
		rescart,
	)
}

func (r *cartroutes) deleteCart(ctx *gin.Context) {
	var (
		function = "delete cart"
		entity   = "cart"

		rid  uuid.UUID
		role consttypes.UserRole
		ok   bool
	)

	muuid, err := uuid.Parse(ctx.Param("muuid"))
	if err != nil {
		utresponse.GeneralInputRequiredError(
			function,
			ctx,
			err,
		)
		return
	}

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

	rid, role, ok = utrole.CartRoleCheck(*roleres)
	if !ok {
		utresponse.GeneralUnauthorized(
			ctx,
			consttypes.ErrUserInvalidRole,
		)
		return
	}

	if role == consttypes.UR_CAREGIVER {
		m, err := r.smmbr.FindByCaregiverID(rid)
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

		rid = m.ID
		role = m.User.Role
	}

	cart, err := r.scart.GetCartByMealIDAndReference(muuid, rid, role)
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

	err = r.scart.Delete(cart.ID)
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

	utresponse.GeneralSuccessDelete(
		entity,
		ctx,
		nil,
	)
}
