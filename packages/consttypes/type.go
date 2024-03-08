package consttypes

import "encoding/json"

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

func (enum ImageType) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}

func (enum PatronType) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}

func (enum OrganizationType) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}

func (enum ResponseStatusType) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}
