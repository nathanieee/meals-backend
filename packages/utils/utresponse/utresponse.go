package utresponse

import (
	"errors"
	"fmt"
	"net/http"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/* -------------------------------------------------------------------------- */
/*                              success responses                             */
/* -------------------------------------------------------------------------- */

type (
	SuccessRes struct {
		Status  consttypes.ResponseStatusType `json:"status"`
		Message string                        `json:"message"`
		Data    any                           `json:"data,omitempty"`
		Header  uttoken.TokenHeader           `json:"-"`
	}
)

func SuccessResponse(ctx *gin.Context, code int, res SuccessRes) {
	if res.Header != (uttoken.TokenHeader{}) {
		ctx.Request.Header.Set(consttypes.T_REFRESH, res.Header.RefreshToken)
		ctx.Request.Header.Set("refresh-token-expired", res.Header.RefreshTokenExpires.String())
		ctx.Request.Header.Set(consttypes.T_ACCESS, "Bearer "+res.Header.AccessToken)
		ctx.Request.Header.Set("expired-at", res.Header.AccessTokenExpires.String())
	}
	ctx.JSON(code, res)
}

func GeneralSuccess(
	function string,
	ctx *gin.Context,
	data any,
) {
	SuccessResponse(ctx, http.StatusOK, SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: fmt.Sprintf("Success %s", function),
		Data:    data,
	})
}

func GeneralSuccessCreate(
	entity string,
	ctx *gin.Context,
	data any,
) {
	SuccessResponse(ctx, http.StatusCreated, SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: fmt.Sprintf("Success creating a new %s", entity),
		Data:    data,
	})
}

func GeneralSuccessUpdate(
	entity string,
	ctx *gin.Context,
	data any,
) {
	SuccessResponse(ctx, http.StatusOK, SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: fmt.Sprintf("Success updating %s", entity),
		Data:    data,
	})
}

func GeneralSuccessFetch(
	entity string,
	ctx *gin.Context,
	data any,
) {
	SuccessResponse(ctx, http.StatusOK, SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: fmt.Sprintf("Success fetching %s", entity),
		Data:    data,
	})
}

func GeneralSuccessDelete(
	entity string,
	ctx *gin.Context,
	data any,
) {
	SuccessResponse(ctx, http.StatusOK, SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: fmt.Sprintf("Success deleting %s", entity),
		Data:    data,
	})
}

func GeneralSuccessAuth(
	function string,
	ctx *gin.Context,
	data any,
	header *uttoken.TokenHeader,
) {
	SuccessResponse(ctx, http.StatusOK, SuccessRes{
		Status:  consttypes.RST_SUCCESS,
		Message: fmt.Sprintf("Success on %s", function),
		Data:    data,
		Header:  *header,
	})
}

/* -------------------------------------------------------------------------- */
/*                               error responses                              */
/* -------------------------------------------------------------------------- */

type (
	ValidationErrorMessage struct {
		Namespace string `json:"namespace"`
		Field     string `json:"field"`
		Message   string `json:"message"`
	}

	ErrorRes struct {
		Status  consttypes.ResponseStatusType `json:"status"`
		Message string                        `json:"message"`
		Data    *ErrorData                    `json:"data,omitempty"`
	}

	ErrorData struct {
		Debug  *error `json:"debug,omitempty"`
		Errors any    `json:"errors"`
	}
)

func ErrorResponse(ctx *gin.Context, code int, res ErrorRes) {
	ctx.JSON(code, res)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "should be a valid email address"
	case "file":
		return "should be a valid file"
	case "uuid":
		return "should be a valid UUID"
	case "lte":
		return "should be less than " + fe.Param()
	case "gte":
		return "should be greater than " + fe.Param()
	case "len":
		return "should be " + fe.Param() + " character(s) long"
	case "eqfield":
		return "should be equal to " + fe.Param()
	}
	return "unknown error"
}

func ValidationResponse(err error) []ValidationErrorMessage {
	var (
		ve validator.ValidationErrors
	)

	if errors.As(err, &ve) {
		out := make([]ValidationErrorMessage, len(ve))
		for i, fe := range ve {
			out[i] = ValidationErrorMessage{fe.Namespace(), fmt.Sprintf("%s", fe.Field()), getErrorMsg(fe)}
			utlogger.Error(err)
		}
		return out
	}

	return nil
}

func GeneralInputRequiredError(
	function string,
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusBadRequest, ErrorRes{
		Status:  consttypes.RST_ERROR,
		Message: fmt.Sprintf("Input required on %s", function),
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}

func GeneralInternalServerError(
	function string,
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusInternalServerError, ErrorRes{
		Status:  consttypes.RST_ERROR,
		Message: fmt.Sprintf("Something went wrong on %s", function),
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}

func GeneralInvalidRequest(
	function string,
	ctx *gin.Context,
	ve []ValidationErrorMessage,
	err error,
) {
	ErrorResponse(ctx, http.StatusBadRequest, ErrorRes{
		Status:  consttypes.RST_FAIL,
		Message: fmt.Sprintf("Invalid request on %s", function),
		Data: &ErrorData{
			Debug:  &err,
			Errors: ve,
		},
	})
}

func GeneralNotFound(
	entity string,
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusNotFound, ErrorRes{
		Status:  consttypes.RST_FAIL,
		Message: fmt.Sprintf("Entity %s is not found", entity),
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}

func GeneralUnauthorized(
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusUnauthorized, ErrorRes{
		Status:  consttypes.RST_FAIL,
		Message: "You are unauthorized to perform this action",
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}

func GeneralFailedCreate(
	entity string,
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusUnprocessableEntity, ErrorRes{
		Status:  consttypes.RST_FAIL,
		Message: fmt.Sprintf("Failed to create %s", entity),
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}

func GeneralFailedUpdate(
	entity string,
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusUnprocessableEntity, ErrorRes{
		Status:  consttypes.RST_FAIL,
		Message: fmt.Sprintf("Failed to update %s", entity),
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}

func GeneralDuplicate(
	field string,
	ctx *gin.Context,
	err error,
) {
	ErrorResponse(ctx, http.StatusConflict, ErrorRes{
		Status:  consttypes.RST_FAIL,
		Message: fmt.Sprintf("Failed to process, there is a duplicate %s", field),
		Data: &ErrorData{
			Debug:  &err,
			Errors: err.Error(),
		},
	})
}
