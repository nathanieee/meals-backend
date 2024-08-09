package consttypes

type (
	PatronType string
)

const (
	PT_ORGANIZATION PatronType = "Organization"
	PT_PERSONAL     PatronType = "Personal"
)

func (enum PatronType) String() string {
	return string(enum)
}
