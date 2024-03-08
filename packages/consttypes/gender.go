package consttypes

import "encoding/json"

type (
	Gender string
)

const (
	G_MALE   Gender = "Male"
	G_FEMALE Gender = "Female"
	G_OTHER  Gender = "Other"
)

func (enum Gender) String() string {
	jsondata, _ := json.Marshal(enum)
	return string(jsondata)
}
