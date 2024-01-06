package utils

import (
	"errors"
	"fmt"
	"net/http"
	"project-skbackend/packages/custom"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/* -------------------------------------------------------------------------- */
/*                              success responses                             */
/* -------------------------------------------------------------------------- */

type (
	TokenHeader struct {
		AuthToken           string
		AuthTokenExpires    time.Time
		RefreshToken        string
		RefreshTokenExpires time.Time
	}

	SuccessRes struct {
		Message string      `json:"message"`
		Data    any         `json:"data,omitempty"`
		Header  TokenHeader `json:"-"`
	}
)

func SuccessResponse(c *gin.Context, code int, res SuccessRes) {
	if res.Header != (TokenHeader{}) {
		c.Header("refresh-token", res.Header.RefreshToken)
		c.Header("refresh-token-expired", res.Header.RefreshTokenExpires.String())
		c.Header("Authorization", "Bearer "+res.Header.AuthToken)
		c.Header("expired-at", res.Header.AuthTokenExpires.String())
	}
	c.JSON(code, res)
}

/* -------------------------------------------------------------------------- */
/*                               error responses                              */
/* -------------------------------------------------------------------------- */

type (
	ValidationErrorMsg struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	ErrorRes struct {
		Message string `json:"message"`
		Debug   error  `json:"debug,omitempty"`
		Errors  any    `json:"errors"`
	}
)

var (
	// General
	ErrConvertFailed = errors.New("data type conversion failed")

	// Error Field
	ErrFieldIsEmpty             = errors.New("field should not be empty")
	ErrFieldInvalidFormat       = errors.New("field format is invalid")
	ErrFieldInvalidEmailAddress = errors.New("invalid email address format")

	// Token
	ErrTokenExpired      = errors.New("token is expired")
	ErrTokenUnverifiable = errors.New("token is unverifiable")
	ErrTokenMismatch     = errors.New("token is mismatch")
	ErrTokenIsNotTheSame = errors.New("this token is not the same")

	// User
	ErrUserNotFound         = errors.New("user not found")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrUserIDNotFound       = errors.New("unable to assert user ID")
	ErrUserAlreadyExist     = errors.New("user already exists")
	ErrUserAlreadyConfirmed = errors.New("this user is already confirmed")

	// Email
	ErrSendEmailResetRequest        = errors.New("you already requested a reset password email in less than 5 minutes")
	ErrSendEmailVerificationRequest = errors.New("you already requested a verification message in less than 5 minutes")
)

func ErrorResponse(c *gin.Context, code int, res ErrorRes) {
	c.JSON(code, res)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "should be a valid email address"
	case "file":
		return "should be a valid file"
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

func ValidationResponse(err error) []ValidationErrorMsg {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]ValidationErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ValidationErrorMsg{fmt.Sprintf("%s", fe.Field()), getErrorMsg(fe)}
		}
		return out
	}

	return nil
}

func GeneralInputRequiredError(m custom.V_STRING, c *gin.Context, e any) {
	ErrorResponse(c, http.StatusBadRequest, ErrorRes{
		Message: string(m.SuffixSpaceCheck()) + "(input required)",
		Debug:   nil,
		Errors:  e,
	})
}

func GeneralInternalServerError(m custom.V_STRING, c *gin.Context, e any) {
	ErrorResponse(c, http.StatusInternalServerError, ErrorRes{
		Message: string(m.SuffixSpaceCheck()) + "(internal server error)",
		Debug:   nil,
		Errors:  e,
	})
}
