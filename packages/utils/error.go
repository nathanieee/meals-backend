package utils

import "errors"

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
