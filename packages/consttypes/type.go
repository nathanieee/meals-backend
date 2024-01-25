package consttypes

type (
	ImageType          string
	PatronType         string
	OrganizationType   string
	ResponseStatusType string
)

const (
	IT_PROFILE ImageType = "Profile"
	IT_MEAL    ImageType = "Meal"

	PT_ORGANIZATION PatronType = "Organization"
	PT_PERSONAL     PatronType = "Personal"

	OT_NURSINGHOME OrganizationType = "Nursing Home"

	RST_SUCCESS ResponseStatusType = "success"
	RST_FAIL    ResponseStatusType = "fail"
	RST_ERROR   ResponseStatusType = "error"
)
