package consttypes

type (
	Level uint
	Role  string
)

const (

	/* ------------------------------- user level ------------------------------- */

	ADMIN Level = 1
	USER  Level = 2

	/* -------------------------------- user role ------------------------------- */

	CAREGIVER     Role = "Caregiver"
	MEMBER        Role = "Member"
	PARTNER       Role = "Partner"
	PATRON        Role = "Patron"
	ADMINISTRATOR Role = "Administrator"
)
