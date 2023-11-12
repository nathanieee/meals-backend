package consttypes

type (
	ImageType        string
	PatronType       string
	OrganizationType string
)

const (
	IT_PROFILE ImageType = "Profile"
	IT_MEAL    ImageType = "Meal"

	PT_ORGANIZATION PatronType = "Organization"
	PT_PERSONAL     PatronType = "Personal"

	OT_NURSINGHOME OrganizationType = "Nursing Home"
)
