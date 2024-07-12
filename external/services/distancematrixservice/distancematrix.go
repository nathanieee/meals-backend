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
	url := fmt.Sprintf("%sgeocode/json?address=%s, %s&key=%s", s.url, loc.Lat, loc.Lng, s.apikey)

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

	utlogger.Info(url)

	if resp.StatusCode != http.StatusOK {
		return nil, consttypes.ErrUnexpectedStatusCode(resp.StatusCode)
	}

	var res *responses.Geocode
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
