package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"project-skbackend/packages/consttypes"
	"sync"
	"time"
)

type (
	TelegramService struct {
		apikey   string
		url      string
		tochatid string

		httpclient *http.Client
	}

	ITelegramService interface {
		SendMessage(msg string) error
	}
)

func NewTelegramService() *TelegramService {
	return &TelegramService{
		apikey:   "bot7046722853:AAGvIquvtNmR8Ttaq7kfZnDb9Zn5D6dQui4",
		url:      "https://api.telegram.org",
		tochatid: "6390863276",

		httpclient: &http.Client{
			Timeout: time.Second * time.Duration(30), // Example: Timeout after 10 seconds
		},
	}
}

func (s *TelegramService) SendMessage(msg string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	encodedMsg := url.QueryEscape(msg)
	url := fmt.Sprintf(
		"%s/%s/sendMessage?chat_id=%s&text=%s",
		s.url, s.apikey, s.tochatid, encodedMsg,
	)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		errChan <- consttypes.ErrFailedToDeclareNewRequest
		return
	}

	resp, err := s.httpclient.Do(req)
	if err != nil {
		errChan <- consttypes.ErrFailedToCallExternalAPI
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- consttypes.ErrUnexpectedStatusCode(resp.StatusCode)
		return
	}

	errChan <- nil
}
