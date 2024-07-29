package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"project-skbackend/packages/consttypes"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
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
	godotenv.Load()

	timeout := os.Getenv("TG_TIMEOUT")
	timeoutint, _ := strconv.Atoi(timeout)

	return &TelegramService{
		apikey:   os.Getenv("TG_API_KEY"),
		url:      os.Getenv("TG_BASE_URL"),
		tochatid: os.Getenv("TG_TO_CHAT_ID"),

		httpclient: &http.Client{
			Timeout: time.Second * time.Duration(timeoutint), // Example: Timeout after 10 seconds
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
