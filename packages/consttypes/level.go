package consttypes

type (
	UserRole uint
)

const (

	/* -------------------------------- user role ------------------------------- */

	UL_USER          UserRole = 0
	UL_ADMINISTRATOR UserRole = 1
	UL_CAREGIVER     UserRole = 2
	UL_MEMBER        UserRole = 3
	UL_PARTNER       UserRole = 4
	UL_PATRON        UserRole = 5
)
