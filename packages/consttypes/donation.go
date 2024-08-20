package consttypes

type (
	DonationStatus string
)

const (
	DS_ACCEPTED DonationStatus = "Accepted"
	DS_PENDING  DonationStatus = "Pending"
	DS_REJECTED DonationStatus = "Rejected"
)

func (enum DonationStatus) String() string {
	return string(enum)
}
