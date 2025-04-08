package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type webhookSender struct {
	httpClient    HTTPClient
	url           string
	verbose       bool
	notifySuccess bool
}

func NewWebhookSender(httpClient HTTPClient, url string, verbose bool, notifySuccess bool) (*webhookSender, error) {
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

func (s *webhookSender) send(m *Results, failedCount int) error {
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

	bodyReader := bytes.NewReader(jsonData)

	if s.verbose {
		log.Printf("prepared webhook data: %s", string(jsonData))
	}

	err = handleSubmit(s.httpClient, http.MethodPost, s.url, bodyReader)
	if err != nil {
		return fmt.Errorf("error pushing results to webhook: %w", err)
	}

	log.Println("results pushed to webhook")

	return nil
}
