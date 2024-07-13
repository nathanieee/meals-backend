package responses

type (
	Geolocation struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	}

	DistanceMatrix struct {
		OriginAddresses      string `json:"origin_addresses"`
		DestinationAddresses string `json:"destination_addresses"`

		Distance Distance `json:"distance"`

		Duration Duration `json:"duration"`
	}

	Distance struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	}

	Duration struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	}
)
