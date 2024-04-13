package consttypes

import (
	"fmt"
	"os"
)

var (
	// variable
	resetPasswordCooldown = os.Getenv("API_RESET_PASSWORD_COOLDOWN")

	// general
	ErrConvertFailed    = fmt.Errorf("data type conversion failed")
	ErrInvalidReference = fmt.Errorf("invalid reference")

	// field
	ErrFieldIsEmpty             = fmt.Errorf("field should not be empty")
	ErrFieldInvalidFormat       = fmt.Errorf("field format is invalid")
	ErrFieldInvalidEmailAddress = fmt.Errorf("invalid email address format")

	// token
	ErrTokenExpired               = fmt.Errorf("token is expired")
	ErrTokenUnverifiable          = fmt.Errorf("token is unverifiable")
	ErrTokenMismatch              = fmt.Errorf("token is mismatch")
	ErrTokenNotFound              = fmt.Errorf("token is not found")
	ErrTokenInvalidFormat         = fmt.Errorf("token format is invalid")
	ErrTokenCannotDecodePublicKey = fmt.Errorf("cannot decode token public key")

	// user
	ErrUserNotFound         = fmt.Errorf("user not found")
	ErrIncorrectPassword    = fmt.Errorf("incorrect password")
	ErrUserIDNotFound       = fmt.Errorf("user ID is not found")
	ErrUserAlreadyExist     = fmt.Errorf("user already exists")
	ErrUserAlreadyConfirmed = fmt.Errorf("this user is already confirmed")
	ErrUserNotSignedIn      = fmt.Errorf("you are not signed in")
	ErrUserInvalidRole      = fmt.Errorf("invalid user role")

	// file
	ErrInvalidFileType = fmt.Errorf("invalid file type")

	// email
	ErrCannotChangeEmail = fmt.Errorf("cannot change existing email")
	ErrTooQuickSendEmail = fmt.Errorf("an email was sent just under %v minutes ago", resetPasswordCooldown)
)
