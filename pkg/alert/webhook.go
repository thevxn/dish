package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type webhookSender struct {
	httpClient *http.Client
	url        string
	verbose    bool
	failedOnly bool
}

func NewWebhookSender(httpClient *http.Client, url string, verbose bool, failedOnly bool) *webhookSender {
	return &webhookSender{
		httpClient,
		url,
		verbose,
		failedOnly,
	}
}

func (s *webhookSender) send(m Results, failedCount int) error {
	// If there are no failed sockets and we only wish to be notified when they fail, there is nothing to do
	if failedCount == 0 && s.failedOnly {
		log.Printf("%T: no failed sockets and failedOnly == true, nothing will be sent", s)
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if s.verbose {
		log.Printf("prepared webhook data: %v", string(jsonData))
	}

	res, err := s.httpClient.Post(s.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code received from webhook (expected: %d, got: %d)", http.StatusOK, res.StatusCode)
	}

	// Write the body to console if verbose flag set
	if s.verbose {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		log.Println("webhook response:", string(body))
	}

	log.Println("results pushed to webhook")

	return nil
}
