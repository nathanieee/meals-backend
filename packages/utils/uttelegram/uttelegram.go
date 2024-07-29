package uttelegram

import (
	"fmt"
	"project-skbackend/external/services/telegram"
	"sync"
)

var (
	stg = telegram.NewTelegramService()
)

func SendMessage(msg string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 10) // Buffer size to handle errors from multiple goroutines

	wg.Add(1)
	go stg.SendMessage(msg, &wg, errChan)

	wg.Wait()
	close(errChan)

	// Handle errors
	for err := range errChan {
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	return nil
}
