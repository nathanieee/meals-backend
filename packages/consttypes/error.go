package consttypes

import (
	"fmt"
	"os"
)

var (
	// variable
	resetPasswordCooldown = os.Getenv("API_RESET_PASSWORD_COOLDOWN")

	// queues
	ErrFailedToPublishMessage = fmt.Errorf("failed to publish message")

	// general
	ErrConvertFailed          = fmt.Errorf("data type conversion failed")
	ErrInvalidReference       = fmt.Errorf("invalid reference")
	ErrUnauthorized           = fmt.Errorf("you are unauthorized to access this resource")
	ErrAccountIsNotVerified   = fmt.Errorf("your account is not verified yet")
	ErrInvalidEmailOrPassword = fmt.Errorf("invalid email or password")
	ErrFailedToGetUserName    = fmt.Errorf("failed to get user's name")

	// field
	ErrFieldIsEmpty             = fmt.Errorf("field should not be empty")
	ErrFieldInvalidFormat       = fmt.Errorf("field format is invalid")
	ErrFieldInvalidEmailAddress = fmt.Errorf("invalid email address format")

	// token
	ErrTokenExpired               = fmt.Errorf("token is expired")
	ErrTokenUnverifiable          = fmt.Errorf("token is unverifiable")
	ErrTokenMismatch              = fmt.Errorf("token is mismatch")
	ErrTokenIsNotTheSame          = fmt.Errorf("token is not the same")
	ErrTokenIsExpired             = fmt.Errorf("token is expired")
	ErrTokenNotFound              = fmt.Errorf("token is not found")
	ErrTokenInvalidFormat         = fmt.Errorf("token format is invalid")
	ErrTokenCannotDecodePublicKey = fmt.Errorf("cannot decode token public key")
	ErrFailedToGenerateToken      = fmt.Errorf("failed to generate token")

	// orders
	ErrFailedToCreateOrder   = fmt.Errorf("failed to create order")
	ErrFailedToReadOrder     = fmt.Errorf("failed to read orders")
	ErrOrderNotFound         = fmt.Errorf("order not found")
	ErrFailedToDeleteOrder   = fmt.Errorf("failed to delete order")
	ErrFailedToFindAllOrders = fmt.Errorf("failed to find all orders")

	// partners
	ErrPartnerNotFound = fmt.Errorf("partner not found")

	// members
	ErrMemberNotFound         = fmt.Errorf("member not found")
	ErrFailedToCreateMember   = fmt.Errorf("failed to create member")
	ErrFailedToReadMembers    = fmt.Errorf("failed to read members")
	ErrFailedToUpdateMember   = fmt.Errorf("failed to update member")
	ErrFailedToDeleteMember   = fmt.Errorf("failed to delete member")
	ErrFailedToFindAllMembers = fmt.Errorf("failed to find all members")

	// caregivers
	ErrCaregiverNotFound = fmt.Errorf("caregiver not found")

	// meals
	ErrMealsNotFound        = fmt.Errorf("meals not found")
	ErrFailedToCreateMeal   = fmt.Errorf("failed to create meal")
	ErrFailedToReadMeals    = fmt.Errorf("failed to read meals")
	ErrFailedToUpdateMeal   = fmt.Errorf("failed to update meal")
	ErrFailedToDeleteMeal   = fmt.Errorf("failed to delete meal")
	ErrFailedToFindAllMeals = fmt.Errorf("failed to find all meals")

	// illnesses
	ErrIllnessNotFound = fmt.Errorf("illness not found")

	// allergies
	ErrAllergiesNotFound = fmt.Errorf("allergies not found")

	// carts
	ErrGettingCart        = fmt.Errorf("failed to get cart")
	ErrFailedToUpdateCart = fmt.Errorf("failed to update cart")
	ErrFailedToCreateCart = fmt.Errorf("failed to create cart")
	ErrFailedToReadCart   = fmt.Errorf("failed to read cart")
	ErrCartNotFound       = fmt.Errorf("cart not found")
	ErrFailedToDeleteCart = fmt.Errorf("failed to delete cart")

	// organizations
	ErrOrganizationNotFound = fmt.Errorf("organization not found")

	// user
	ErrUserNotFound         = fmt.Errorf("user not found")
	ErrIncorrectPassword    = fmt.Errorf("incorrect password")
	ErrUserIDNotFound       = fmt.Errorf("user ID is not found")
	ErrUserAlreadyExist     = fmt.Errorf("user already exists")
	ErrUserAlreadyConfirmed = fmt.Errorf("this user is already confirmed")
	ErrUserNotSignedIn      = fmt.Errorf("you are not signed in")
	ErrUserInvalidRole      = fmt.Errorf("invalid user role")
	ErrFailedToUpdateUser   = fmt.Errorf("failed to update user")

	// file
	ErrInvalidFileType         = fmt.Errorf("invalid file type")
	ErrFailedToUploadFile      = fmt.Errorf("failed to upload file")
	ErrFailedToCreateDirectory = fmt.Errorf("failed to create directory")
	ErrFailedToParseFile       = fmt.Errorf("failed to parse file")
	ErrFailedToWriteFile       = fmt.Errorf("failed to write file")

	// caches
	ErrFailedToSetCache = fmt.Errorf("failed to set cache")
	ErrFailedToGetCache = fmt.Errorf("failed to get cache")

	// email
	ErrCannotChangeEmail = fmt.Errorf("cannot change existing email")
	ErrTooQuickSendEmail = fmt.Errorf("an email was sent just under %v minutes ago", resetPasswordCooldown)
	ErrDuplicateEmail    = fmt.Errorf("email address already exists")
	ErrFailedToSendEmail = fmt.Errorf("failed to send email")
)
