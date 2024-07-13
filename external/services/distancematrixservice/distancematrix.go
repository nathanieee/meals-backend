package distancematrixservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-skbackend/configs"
	"project-skbackend/external/controllers/requests"
	"project-skbackend/external/controllers/responses"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"time"
)

type (
	DistanceMatrixService struct {
		apikey string
		url    string

		httpclient *http.Client
	}

	IDistanceMatrixService interface {
		GetGeocoding(loc requests.Geolocation) (*responses.Geocode, error)
		GetDistanceMatrix(dismat requests.DistanceMatrix) (*responses.DistanceMatrix, error)
	}
)

func NewDistanceMatrixService(
	cfg *configs.Config,
) *DistanceMatrixService {
	return &DistanceMatrixService{
		apikey: cfg.DistanceMatrix.APIKey,
		url:    cfg.DistanceMatrix.BaseURL,

		httpclient: &http.Client{
			Timeout: time.Second * time.Duration(cfg.DistanceMatrix.Timeout), // Example: Timeout after 10 seconds
		},
	}
}

func (s *DistanceMatrixService) GetGeocoding(loc requests.Geolocation) (*responses.Geocode, error) {
	url := fmt.Sprintf(
		"%sgeocode/json?address=%s, %s&key=%s",
		s.url, loc.Lat, loc.Lng, s.apikey,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utlogger.Error(err)
		return nil, consttypes.ErrFailedToDeclareNewRequest
	}

	req.Header.Set("Authorization", "Bearer "+s.apikey)
	resp, err := s.httpclient.Do(req)
	if err != nil {
		utlogger.Error(err)
		return nil, consttypes.ErrFailedToCallExternalAPI
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, consttypes.ErrUnexpectedStatusCode(resp.StatusCode)
	}

	var res *responses.Geocode
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * return nil if the status is not OK
	if res.Status != consttypes.DMS_OK.String() {
		utlogger.Info("invalid response: ", res)
		return nil, consttypes.ErrInvalidGeolocation
	}

	return res, nil
}

func (s *DistanceMatrixService) GetDistanceMatrix(dismat requests.DistanceMatrix) (*responses.DistanceMatrix, error) {
	url := fmt.Sprintf(
		"%sdistancematrix/json?origins=%s, %s&destinations=%s, %s&key=%s",
		s.url, dismat.Origins.Lat, dismat.Origins.Lng, dismat.Destinations.Lat, dismat.Destinations.Lng, s.apikey,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utlogger.Error(err)
		return nil, consttypes.ErrFailedToDeclareNewRequest
	}

	req.Header.Set("Authorization", "Bearer "+s.apikey)
	resp, err := s.httpclient.Do(req)
	if err != nil {
		utlogger.Error(err)
		return nil, consttypes.ErrFailedToCallExternalAPI
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, consttypes.ErrUnexpectedStatusCode(resp.StatusCode)
	}

	var res *responses.DistanceMatrix
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * return nil if the status is not OK
	if res.Rows[0].Elements[0].Status != consttypes.DMS_OK.String() || res.Status != consttypes.DMS_OK.String() {
		utlogger.Info("invalid response: ", res)
		return nil, consttypes.ErrInvalidDistanceMatrix
	}

	return res, nil
}
