package consttypes

type (
	Gender string
)

const (
	G_MALE   Gender = "Male"
	G_FEMALE Gender = "Female"
	G_OTHER  Gender = "Other"
)

func (enum Gender) String() string {
	return string(enum)
}
