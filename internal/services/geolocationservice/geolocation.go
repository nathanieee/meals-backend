package geolocationservice

import (
	"project-skbackend/external/services/distancematrixservice"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
)

type (
	GeolocationService struct {
		sdsmx distancematrixservice.IDistanceMatrixService
	}

	IGeolocationService interface {
		GetGeolocation(loc requests.Geolocation) (*responses.AddressDetail, error)
	}
)

func NewGeolocationService(
	sdsmx distancematrixservice.IDistanceMatrixService,
) *GeolocationService {
	return &GeolocationService{
		sdsmx: sdsmx,
	}
}

func (s *GeolocationService) GetGeolocation(loc requests.Geolocation) (*responses.AddressDetail, error) {
	geoloc := loc.ToDistanceMatrixGeolocation()

	geocode, err := s.sdsmx.GetGeocoding(*geoloc)
	if err != nil {
		return nil, err
	}

	address := geocode.ToAddressDetail()

	return address, nil
}
