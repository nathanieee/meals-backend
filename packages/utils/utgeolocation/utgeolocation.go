package utgeolocation

import (
	"project-skbackend/configs"
	"project-skbackend/external/controllers/exrequests"
	"project-skbackend/external/controllers/exresponses"
	"project-skbackend/external/services/distancematrixservice"
	"project-skbackend/packages/consttypes"
)

var (
	cfg = configs.GetInstance()

	sdsmt = distancematrixservice.NewDistanceMatrixService(cfg)
)

func GetGeolocation(geoloc exrequests.Geolocation) (*exresponses.AddressDetail, error) {
	geocode, err := sdsmt.GetGeocoding(geoloc)
	if err != nil {
		return nil, err
	}

	if geocode == nil {
		return nil, consttypes.ErrGeolocationNotFound
	}

	address := geocode.ToAddressDetail()

	return address, nil
}

func GetLocationDistance(dismat exrequests.DistanceMatrix) (*exresponses.RouteDetails, error) {
	distancematrix, err := sdsmt.GetDistanceMatrix(dismat)
	if err != nil {
		return nil, err
	}

	if distancematrix == nil {
		return nil, consttypes.ErrInvalidDistanceMatrix
	}

	routedetails := distancematrix.ToRouteDetails()

	return routedetails, nil
}
