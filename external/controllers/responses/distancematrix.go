package responses

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
)

type (
	/* -------------------------------------------------------------------------- */
	/*                                   geocode                                  */
	/* -------------------------------------------------------------------------- */
	Geocode struct {
		Results []GeocodeResult `json:"result"`
		Status  string          `json:"status"`
	}

	GeocodeResult struct {
		GeocodeAddressComponents []*GeocodeAddressComponent `json:"address_components"`
		FormattedAddress         string                     `json:"formatted_address"`
		Geometry                 *GeocodeGeometry           `json:"geometry"`
		Types                    []string                   `json:"types"`
	}

	GeocodeGeometry struct {
		Location     *Geolocation `json:"location"`
		LocationType string       `json:"location_type"`
	}

	GeocodeViewport struct {
		Northeast *Geolocation `json:"northeast"`
		Southwest *Geolocation `json:"southwest"`
	}

	Geolocation struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	GeocodeAddressComponent struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	}

	/* -------------------------------------------------------------------------- */
	/*                               distance matrix                              */
	/* -------------------------------------------------------------------------- */
	DistanceMatrix struct {
		OriginAddresses      []string            `json:"origin_addresses"`
		DestinationAddresses []string            `json:"destination_addresses"`
		Rows                 []DistanceMatrixRow `json:"rows"`
		Status               string              `json:"status"`
	}

	DistanceMatrixRow struct {
		Elements []DistanceMatrixElement `json:"elements"`
	}

	DistanceMatrixElement struct {
		Distance    *DistanceMatrixDistance `json:"distance"`
		Duration    *DistanceMatrixDuration `json:"duration"`
		Origin      string                  `json:"origin"`
		Destination string                  `json:"destination"`
		Status      string                  `json:"status"`
	}

	DistanceMatrixDistance struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	}

	DistanceMatrixDuration struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	}
)

func (res *Geocode) ToAddressDetail() *responses.AddressDetail {
	var (
		postcode = ""
		country  = ""
		result   = res.Results[0]
	)

	for _, component := range result.GeocodeAddressComponents {
		if component.Types[0] == "postcode" {
			postcode = component.LongName
		}

		if component.Types[0] == "country" {
			country = component.LongName
		}
	}

	return &responses.AddressDetail{
		Geolocation: responses.Geolocation{
			Longitude: fmt.Sprintf("%.6f", result.Geometry.Location.Lng),
			Latitude:  fmt.Sprintf("%.6f", result.Geometry.Location.Lat),
		},
		FormattedAddress: result.FormattedAddress,
		PostCode:         postcode,
		Country:          country,
	}
}
