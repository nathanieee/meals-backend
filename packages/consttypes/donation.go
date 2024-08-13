package consttypes

type (
	DonationStatus string
)

const (
	DS_ACCEPTED DonationStatus = "Accepted"
	DS_REJECTED DonationStatus = "Rejected"
)

func (enum DonationStatus) String() string {
	return string(enum)
}
