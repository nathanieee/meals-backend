package consttypes

type (
	UserRole uint
)

const (

	/* -------------------------------- user role ------------------------------- */

	UR_USER          UserRole = 0
	UR_ADMINISTRATOR UserRole = 1
	UR_CAREGIVER     UserRole = 2
	UR_MEMBER        UserRole = 3
	UR_PARTNER       UserRole = 4
	UR_PATRON        UserRole = 5
)
