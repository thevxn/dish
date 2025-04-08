package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiSender struct {
	httpClient    HTTPClient
	url           string
	headerName    string
	headerValue   string
	verbose       bool
	notifySuccess bool
}

func NewAPISender(httpClient HTTPClient, url string, headerName string, headerValue string, verbose bool, notifySuccess bool) (*apiSender, error) {
	parsedURL, err := parseAndValidateURL(url, nil)
	if err != nil {
		return nil, err
	}

	return &apiSender{
		httpClient:    httpClient,
		url:           parsedURL.String(),
		headerName:    headerName,
		headerValue:   headerValue,
		verbose:       verbose,
		notifySuccess: notifySuccess,
	}, nil
}

func (s *apiSender) send(m *Results, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		if s.verbose {
			log.Println("no sockets failed, nothing will be sent to remote API")
		}
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	bodyReader := bytes.NewReader(jsonData)

	if s.verbose {
		log.Printf("prepared remote API data: %s", string(jsonData))
	}

	// If custom header & value is provided (mostly used for auth purposes), include it in the request
	opts := []func(*submitOptions){}
	if s.headerName != "" && s.headerValue != "" {
		opts = append(opts, withHeader(s.headerName, s.headerValue))
	}

	err = handleSubmit(s.httpClient, http.MethodPost, s.url, bodyReader, opts...)
	if err != nil {
		return fmt.Errorf("error pushing results to remote API: %w", err)
	}

	log.Println("results pushed to remote API")

	return nil
}
