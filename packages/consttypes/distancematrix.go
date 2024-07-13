package consttypes

type (
	DistanceMatrixStatus string
	DistanceMatrixTypes  string
)

const (
	DMS_OK           DistanceMatrixStatus = "OK"
	DMS_ZERO_RESULTS DistanceMatrixStatus = "ZERO_RESULTS"

	DMT_POSTCODE DistanceMatrixTypes = "postcode"
	DMT_COUNTRY  DistanceMatrixTypes = "country"
	DMT_HOUSE    DistanceMatrixTypes = "house"
	DMT_LOCALITY DistanceMatrixTypes = "locality"
)

func (d DistanceMatrixStatus) String() string {
	return string(d)
}

func (d DistanceMatrixTypes) String() string {
	return string(d)
}
