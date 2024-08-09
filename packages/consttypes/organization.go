package consttypes

type (
	OrganizationType string
)

const (
	OT_NURSINGHOME OrganizationType = "Nursing Home"
)

func (enum OrganizationType) String() string {
	return string(enum)
}
