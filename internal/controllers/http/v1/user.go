package controllers

import (
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/middlewares"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utrequest"
	"project-skbackend/packages/utils/utresponse"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	userroutes struct {
		cfg   *configs.Config
		suser userservice.IUserService
		smail mailservice.IMailService
	}
)

func newUserRoutes(
	rg *gin.RouterGroup,
	cfg *configs.Config,
	suser userservice.IUserService,
	smail mailservice.IMailService,
) {
	r := &userroutes{
		cfg:   cfg,
		suser: suser,
		smail: smail,
	}

	gadmn := rg.Group("users")
	{
		gadmn.GET("", r.getUser)
		gadmn.POST("", r.createUser)
	}

	guser := rg.Group("users")
	guser.Use(middlewares.JWTAuthMiddleware(
		cfg,
		consttypes.UR_USER,
	))
	{
		guser.GET("me", r.getCurrentUser)
		guser.DELETE("delete", r.deleteUser)
	}
}

func (r *userroutes) createUser(ctx *gin.Context) {
	var (
		req requests.CreateUser
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ve := utresponse.ValidationResponse(err)
		utresponse.GeneralInvalidRequest(
			"create user",
			ctx,
			ve,
			err,
		)
		return
	}

	ures, err := r.suser.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			utresponse.ErrorResponse(ctx, http.StatusConflict, utresponse.ErrorRes{
				Status:  consttypes.RST_ERROR,
				Message: "Duplicate email",
				Data: &utresponse.ErrorData{
					Debug:  &err,
					Errors: err.Error(),
				},
			})
		} else {
			utresponse.ErrorResponse(ctx, http.StatusInternalServerError, utresponse.ErrorRes{
				Status:  consttypes.RST_ERROR,
				Message: "Something went wrong",
				Data: &utresponse.ErrorData{
					Debug:  &err,
					Errors: err.Error(),
				},
			})
		}
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Creating new user",
		Data:    ures,
	})
}

func (r *userroutes) getUser(ctx *gin.Context) {
	paginationReq := utrequest.GeneratePaginationFromRequest(ctx)

	users, err := r.suser.FindAll(paginationReq)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "users not found",
			Data: &utresponse.ErrorData{
				Debug:  &err,
				Errors: err.Error(),
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Get Users",
		Data:    users,
	})
}

func (r *userroutes) getCurrentUser(ctx *gin.Context) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error getting user",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "user not found",
			},
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.User)
	if !ok {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error getting user",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "unable to assert user id",
			},
		})
		return
	}

	user, err := r.suser.FindByID(loggedInUser.ID)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "User not found",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "user not found",
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Get User",
		Data:    user,
	})
}

func (r *userroutes) deleteUser(ctx *gin.Context) {
	ctxUser, exists := ctx.Get("user")
	if !exists {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error getting user",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "user not found",
			},
		})
		return
	}

	loggedInUser, ok := ctxUser.(responses.User)
	if !ok {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Error getting user",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "unable to assert user id",
			},
		})
		return
	}

	if loggedInUser.Role == consttypes.UR_ADMIN {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Something went wrong",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "unable to assert user id",
			},
		})
		return
	}

	err := r.suser.Delete(loggedInUser.ID)
	if err != nil {
		utresponse.ErrorResponse(ctx, http.StatusNotFound, utresponse.ErrorRes{
			Status:  consttypes.RST_ERROR,
			Message: "Something went wrong",
			Data: &utresponse.ErrorData{
				Debug:  nil,
				Errors: "failed deleting user",
			},
		})
		return
	}

	utresponse.SuccessResponse(ctx, http.StatusOK, utresponse.SuccessRes{
		Message: "Success Delete User",
		Data:    nil,
	})
}
