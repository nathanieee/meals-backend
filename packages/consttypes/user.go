package consttypes

type (
	Level uint
	Role  string
)

const (
	/* ------------------------------- USER LEVEL ------------------------------- */

	ADMIN Level = 1
	USER  Level = 2

	/* -------------------------------- USER ROLE ------------------------------- */

	CAREGIVER     Role = "Caregiver"
	MEMBER        Role = "Member"
	PARTNER       Role = "Partner"
	PATRON        Role = "Patron"
	ADMINISTRATOR Role = "Administrator"
)
