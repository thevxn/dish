package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type apiSender struct {
	httpClient  *http.Client
	url         string
	headerName  string
	headerValue string
	verbose     bool
}

func NewApiSender(httpClient *http.Client, url string, headerName string, headerValue string, verbose bool) *apiSender {
	return &apiSender{
		httpClient,
		url,
		headerName,
		headerValue,
		verbose,
	}
}

func (s *apiSender) send(m Results, failedCount int) error {

	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	bodyReader := bytes.NewReader(jsonData)

	if s.verbose {
		log.Printf("prepared remote API data: %v", string(jsonData))
	}

	url := s.url

	// TODO: move?
	// TODO: also add to PGW, TG and webhooks?
	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return err
	}

	match := regex.MatchString(url)
	if !match {
		// TODO: mention the protocol must be included?
		return fmt.Errorf("invalid remote API URL, results have not been pushed")
	}

	// Push results
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(s.headerName, s.headerValue)

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
