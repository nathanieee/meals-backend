package geolocationservice

import (
	"project-skbackend/external/services/distancematrixservice"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/packages/consttypes"
)

type (
	GeolocationService struct {
		sdsmx distancematrixservice.IDistanceMatrixService
	}

	IGeolocationService interface {
		GetGeolocation(loc requests.Geolocation) (*responses.AddressDetail, error)
		GetLocationDistance(dismat requests.DistanceMatrix) (*responses.DistanceMatrix, error)
	}
)

func NewGeolocationService(
	sdsmx distancematrixservice.IDistanceMatrixService,
) *GeolocationService {
	return &GeolocationService{
		sdsmx: sdsmx,
	}
}

func (s *GeolocationService) GetGeolocation(geoloc requests.Geolocation) (*responses.AddressDetail, error) {
	var (
		geolocnew = geoloc.ToDistanceMatrixGeolocation()
	)

	geocode, err := s.sdsmx.GetGeocoding(*geolocnew)
	if err != nil {
		return nil, err
	}

	if geocode == nil {
		return nil, consttypes.ErrGeolocationNotFound
	}

	address := geocode.ToAddressDetail()

	return address, nil
}

func (s *GeolocationService) GetLocationDistance(dismat requests.DistanceMatrix) (*responses.DistanceMatrix, error) {
	var (
		dismatnew = dismat.ToDistanceMatrix()
	)

	distancematrix, err := s.sdsmx.GetDistanceMatrix(*dismatnew)
	if err != nil {
		return nil, err
	}

	if distancematrix == nil {
		return nil, consttypes.ErrInvalidDistanceMatrix
	}

	distancematrixnew := distancematrix.ToDistanceMatrix()

	return distancematrixnew, nil
}
