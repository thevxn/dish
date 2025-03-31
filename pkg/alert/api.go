package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type apiSender struct {
	httpClient    *http.Client
	url           string
	headerName    string
	headerValue   string
	verbose       bool
	notifySuccess bool
}

func NewApiSender(httpClient *http.Client, url string, headerName string, headerValue string, verbose bool, notifySuccess bool) *apiSender {
	return &apiSender{
		httpClient,
		url,
		headerName,
		headerValue,
		verbose,
		notifySuccess,
	}
}

func (s *apiSender) send(m Results, failedCount int) error {
	// If no checks failed and failedOnly is set to true, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		if s.verbose {
			log.Println("no sockets failed and notifySuccess == false, nothing will be sent to remote API")
		}
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	bodyReader := bytes.NewReader(jsonData)

	if s.verbose {
		log.Printf("prepared remote API data: %v", string(jsonData))
	}

	// Parse and validate the provided remote API url
	parsedURL, err := url.Parse(s.url)
	if err != nil {
		return fmt.Errorf("error parsing remote API url: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("the protocol must be specified in the remote API url (e.g. https://...)")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported protocol for remote API provided: %s", parsedURL.Scheme)
	}

	// Push results
	req, err := http.NewRequest(http.MethodPost, parsedURL.String(), bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// If custom header & value is provided (mostly used for auth purposes), include it in the request
	if s.headerName != "" && s.headerValue != "" {
		req.Header.Set(s.headerName, s.headerValue)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code received from remote API (expected: %d, got: %d)", http.StatusOK, res.StatusCode)
	}

	// Write the body to console if verbose flag set
	if s.verbose {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		log.Println("remote API response:", string(body))
	}

	log.Println("results pushed to remote API")

	return nil
}
