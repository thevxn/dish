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
	httpClient    *http.Client
	url           string
	verbose       bool
	notifySuccess bool
}

func NewWebhookSender(httpClient *http.Client, url string, verbose bool, notifySuccess bool) (*webhookSender, error) {
	parsedURL, err := parseAndValidateURL(url, nil)
	if err != nil {
		return nil, err
	}

	return &webhookSender{
		httpClient:    httpClient,
		url:           parsedURL.String(),
		verbose:       verbose,
		notifySuccess: notifySuccess,
	}, nil
}

func (s *webhookSender) send(m Results, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		if s.verbose {
			log.Printf("no sockets failed, nothing will be sent to webhook")
		}
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if s.verbose {
		log.Printf("prepared webhook data: %s", string(jsonData))
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
